# Copyright (c) 2019 Intel Corporation.

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.


import os
import logging
import cv2
import numpy as np
import json
import threading

from .defect import Defect
from .display_info import DisplayInfo
from libs.base_classifier import BaseClassifier
from openvino.inference_engine import IENetwork, IEPlugin

"""FLANN algorithm classifier
"""

minMatches = 10
D_MISSING = 0
D_SHORT = 1


class Classifier(BaseClassifier):
    """Classifier object
    """

    def __init__(self, classifier_config, input_queue, output_queue):
        """Constructor
        
        Parameters
        ----------
        classifier_config : dict
            Configuration object for the classifier
        input_queue : Queue
            input queue for classifier
        output_queue : Queue
            output queue of classifier
        
        Returns
        -------
            Classification object
        """
        super().__init__(classifier_config, input_queue, output_queue)
        self.log = logging.getLogger('DEFECT_DETECTION')
        self.ref_img = classifier_config["ref_img"]
        self.ref_config_roi = classifier_config["ref_config_roi"]
        self.model_xml = classifier_config["model_xml"]
        self.model_bin = classifier_config["model_bin"]
        self.device = classifier_config["device"]

        # Assert all input parameters exist
        assert os.path.exists(self.ref_img), \
            'Reference image does not exist: {}'.format(self.ref_img)
        assert os.path.exists(self.ref_config_roi), \
            'Reference image ROI file does not exist: {}'.\
            format(self.ref_config_roi)
        assert os.path.exists(self.model_xml), \
            'Tensorflow model missing: {}'.format(self.model_xml)
        assert os.path.exists(self.model_bin), \
            'Tensorflow model bin file missing: {}'.format(self.model_bin)

        # Load ref image
        self.ref_img = cv2.imread(self.ref_img)

        # Load ROI config file
        with open(self.ref_config_roi, 'r') as f:
            self.config_roi = json.load(f)

        # Load OpenVINO model
        self.plugin = IEPlugin(device=self.device.upper(), plugin_dirs="")
        self.net = IENetwork.from_ir(model=self.model_xml, weights=self.model_bin)
        self.input_blob = next(iter(self.net.inputs))
        self.output_blob = next(iter(self.net.outputs))
        self.net.batch_size = 1  # change to enable batch loading
        self.exec_net = self.plugin.load(network=self.net)

        # Initialize keypoint descriptor
        self.brisk = cv2.BRISK_create()
        # Flann arguments
        FLANN_INDEX_LSH = 6
        index_params = dict(algorithm=FLANN_INDEX_LSH,
                            table_number=6,
                            key_size=12,
                            multi_probe_level=1)
        search_params = dict(checks=50)
        self.flann = cv2.FlannBasedMatcher(index_params, search_params)

        # Detect keypoints and extract features reference images
        ref_gray = cv2.cvtColor(self.ref_img, cv2.COLOR_BGR2GRAY)
        self.ref_kp = self.brisk.detect(ref_gray, None)
        self.ref_kp, self.ref_des = self.brisk.compute(ref_gray, self.ref_kp)

    def _calculate_homography(self, frame):
        """ Calculate keypoint matches and homography score between
            selected frame and ref image
        
        Parameters
        ----------
        frame     : Selected frame to compare against ref image
        
        Returns
        ------
        Homography between frame and ref image
        """

        ref_kp = self.ref_kp
        ref_des = self.ref_des
        src_pts = []
        dst_pts = []
        matches = []

        # Perform BRISK on incoming frame
        img_gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        img_kp = self.brisk.detect(img_gray, None)
        img_kp, img_des = self.brisk.compute(img_gray, img_kp)

        # Find matching keypoints
        raw_matches = self.flann.knnMatch(ref_des, img_des, k=2)
        for m in raw_matches:
            # Lowe's Ratio Test
            if len(m) == 2 and m[0].distance < m[1].distance * 0.7:
                matches.append(m[0])

        if len(matches) <= minMatches:
            self.log.debug("""Found less than minimum # of matches
                              required to overlay""")
            return None

        self.log.debug("Number of good matches: {}".format(len(matches)))

        src_pts = np.float32([ref_kp[m.queryIdx].pt for m
                              in matches]).reshape(-1, 1, 2)
        dst_pts = np.float32([img_kp[m.trainIdx].pt for m
                              in matches]).reshape(-1, 1, 2)

        # Calculate homography
        M, mask = cv2.findHomography(src_pts, dst_pts, cv2.RANSAC, 4.0)
        score = float(mask.sum()) / mask.size
        return (M, mask, score)

    def _det_mobilenet(self, defect_roi, ref, test, d_type):
        bndbx = []
        n, c, h, w = self.net.inputs[self.input_blob].shape
        for roi in range(0, len(defect_roi)):
            x, y, x1, y1 = defect_roi[roi]
            test_crop = test[y:y1, x:x1]
            test_img = cv2.resize(test_crop, (w, h))
            test_img = test_img.transpose((2, 0, 1))
            res = self.exec_net.infer(inputs={self.input_blob: test_img})
            res = res[self.output_blob]
            probs = np.squeeze(res)
            if d_type == D_MISSING:
                if probs[(roi * 2) + 1] > 0.7:
                    self.log.debug("******Defects_MISSING" + str(roi + 1))
                    bndbx.append(defect_roi[roi])
            elif d_type == D_SHORT:
                if probs[((roi + 5) * 2) + 1] > 0.7:
                    self.log.debug("******Defects_SHORT" + str(roi + 1))
                    bndbx.append(defect_roi[roi])
        return bndbx

    # Main classification algorithm
    def classify(self):
        """Reads the image frame from input queue for classifier 
        and classifies against the specified reference image.
        """
        thread_id = threading.get_ident()
        self.log.info("Classifier thread ID: {} started...".format(thread_id))

        while not self.stop_event.is_set():
            data = self.input_queue.get()
            self.log.debug("classify data: {}".format(data[:2]))
            metadata = data[1]

            # Convert the buffer into np array.
            np_buffer = np.frombuffer(data[2], dtype=np.uint8)
            encoding = metadata["encoding"]
            if encoding is not None:
                reshape_frame = np.reshape(np_buffer, (np_buffer.shape))
                reshape_frame = cv2.imdecode(reshape_frame, 1)
            else:
                reshape_frame = np.reshape(np_buffer, (int(metadata["height"]),
                                                int(metadata["width"]),
                                                int(metadata["channel"])))

            index = metadata.get("user_data", None) # Angle information of frame

            # used for debugging purposes
            if(index == -1):
                d_info = []
                defects = []
                self.log.info("user_data: {} received...".format(index))
                return d_info, defects
            
            # Read correct reference image to
            # perform keypoint detection and overlay
            ref_img = self.ref_img.copy()
            img_orig = reshape_frame.copy()

            # Calculate homography info between frame and corresponding ref image
            M, mask, score = self._calculate_homography(frame=reshape_frame)

            if score is None:
                self.log.debug("Low homography score. Skipping frame")
                defects = []
                return defects

            h_ref, w_ref, z_ref = ref_img.shape
            pts = np.float32([
                [0, 0],
                [0, h_ref - 1],
                [w_ref - 1, h_ref - 1],
                [w_ref - 1, 0]
                ]).reshape(-1, 1, 2)
            dst = cv2.perspectiveTransform(pts, M)

            # Warp ROI in frame to overlay on ref image
            # Determine top-left, top-right,
            # bottom-left, bottom-right corners to warp ROI
            points = np.int32(dst).reshape(4, 2)
            rect = np.zeros((4, 2), dtype="float32")
            rect[0] = points[0]
            rect[1] = points[3]
            rect[2] = points[2]
            rect[3] = points[1]
            (tl, tr, br, bl) = rect
            maxWidth = w_ref
            maxHeight = h_ref

            # Top down view of ROI
            destination = np.array([
                [0, 0],
                [maxWidth - 1, 0],
                [maxWidth - 1, maxHeight - 1],
                [0, maxHeight - 1]], dtype="float32")
            mat = cv2.getPerspectiveTransform(rect, destination)
            img_warp = cv2.warpPerspective(img_orig, mat, (maxWidth, maxHeight))

            src = ref_img.copy()
            overlay = img_warp.copy()

            # Get defect ROI and threshold values from config file
            missing_ROI = []
            short_ROI = []

            if index == 1:
                # Read ROI of known defect locations :
                # missing components and shorts
                if 'missing' in self.config_roi:
                    missing_ROI = self.config_roi["missing"]["A1_roi"]
                if 'short' in self.config_roi:
                    short_ROI = self.config_roi["short"]["A1_roi"]
            defects = []
            d_info = []

            # Missing component defect classification
            missing_bndbx = []
            if len(missing_ROI) > 0:
                missing_bndbx = self._det_mobilenet(missing_ROI, src,
                                                    overlay, D_MISSING)
                for count in range(0, len(missing_bndbx)):
                    [x, y, x1, y1] = missing_bndbx[count]
                    # Convert coordinates to correspond to un-warped image
                    coord = np.asarray([[[x, y], [x1, y1]]], dtype="float32")
                    new_coord = cv2.perspectiveTransform(coord,
                                                        np.linalg.inv(mat))[0]
                    x = new_coord[0][0]
                    y = new_coord[0][1]
                    x1 = new_coord[1][0]
                    y1 = new_coord[1][1]
                    # Defect class of 0 => BP defect
                    # (x,y) -> top left bounding box coordinates
                    # (x1,y1) -> bottom right bounding box coordinates
                    defects.append(Defect(D_MISSING, (x, y), (x1, y1)))
            # Short defect classification
            short_bndbx = []
            if len(short_ROI) > 0:
                short_bndbx = self._det_mobilenet(short_ROI, src, overlay, D_SHORT)
                for count in range(0, len(short_bndbx)):
                    [x, y, x1, y1] = short_bndbx[count]
                    # Convert coordinates to correspond to un-warped image
                    coord = np.asarray([[[x, y], [x1, y1]]], dtype="float32")
                    new_coord = cv2.perspectiveTransform(coord,
                                                        np.linalg.inv(mat))[0]
                    x = new_coord[0][0]
                    y = new_coord[0][1]
                    x1 = new_coord[1][0]
                    y1 = new_coord[1][1]
                    # Defect class of 1 => MR defect
                    # (x,y) -> top left bounding box coordinates
                    # (x1,y1) -> bottom right bounding box coordinates
                    defects.append(Defect(D_SHORT, (x, y), (x1, y1)))
            
            # If defects exists display HIGH [priority: 2] alert string
            if len(defects) > 0:
                d_info.append(DisplayInfo('DEFECT DETECTED', 2))

            # Set state of random number generator after every frame to overcome
            # probabilistic nature of Flann matcher. This might need to be changed.
            cv2.setRNGSeed(0)

            defects_dict = {}
            defect_res = []
            for d in defects:
                if d.defect_class in defects:
                    defects_dict[d.defect_class] += 1
                else:
                    defects_dict[d.defect_class] = 1
                defect_res.append({
                    'type': d.defect_class,
                    'tl': d.tl,
                    'br': d.br
                })

            display_info = []
            for d in d_info:
                display_info.append({
                    'info': d.info,
                    'priority': d.priority
                })         

            metadata["defects"] = defect_res
            metadata["display_info"] = display_info
            data[1] = metadata
            self.output_queue.put(data)
            self.log.debug("Data: {} added to classifier output queue".format(self.output_queue))

        self.log.info("Classifier thread ID: {} stopped...".format(thread_id))
    