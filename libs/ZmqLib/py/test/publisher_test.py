import cv2
import os
import argparse
import json
import etcd3
from distutils.util import strtobool
from libs.ZmqLib.py.pub.publisher import ZMQ_PUBLISHER


def main():
    """main method"""
    args = parse_args()
    save_to_disk = bool(strtobool(args.save_to_disk))
    etcd = etcd3.client()

    config_list = os.environ[os.environ['pub_topic'].lower()
                             + "_config"].split(',')
    publish_address = ""
    if "tcp" in config_list[0].lower():
        publish_address = "tcp://"+config_list[1]
    elif "ipc" in config_list[0].lower():
        publish_address = "ipc:///"+config_list[1]

    # Initiate the publisher
    publisher = ZMQ_PUBLISHER(publish_address,
                              public_keys_dir=args.public_keys_dir,
                              secret_keys_dir=args.private_keys_dir)

    # Initiate the OpenCV VideoCapture
    if not os.path.exists("./pcb_d2000.avi"):
        print("Video file does not exist...")
    cap = cv2.VideoCapture("./pcb_d2000.avi")
    if not cap.isOpened():
        print("Failed to open video file...")

    iterator = 0
    while True:
        # Read frames from video and publish meta-data and numpy frame
        ret, frame = cap.read()
        if ret:
            metaData = dict(
                dtype=str(frame.dtype),
                shape=frame.shape,
                name="{0}{1}".format('Frame', iterator)
            )
            # Sample etcd tests
            result = json.dumps(metaData)
            etcd.put('metaData', result)
            etcd_metaData = etcd.get('metaData')
            etcd.delete('metaData')
            publish_list = [os.environ['pub_topic'],
                            etcd_metaData[0].decode("utf-8"),
                            frame]

            # Publish the meta-data and numpy frame
            publisher.send(publish_list)

            # Write published image frame to disk
            if save_to_disk:
                outputFilePath = "./{0}".format(metaData['name'])
                with open(outputFilePath, "wb") as outfile:
                    outfile.write(frame)
            iterator += 1
        if not ret:
            cap.release()
            # Looping video to continue publishing
            cap = cv2.VideoCapture("./pcb_d2000.avi")

    # Close socket
    publisher.close_socket()


def parse_args():

    parser = argparse.ArgumentParser()

    parser.add_argument('--public-key', dest='public_keys_dir',
                        default="",
                        help='Public keys directory')

    parser.add_argument('--private-key', dest='private_keys_dir',
                        default="",
                        help='Public keys directory')

    parser.add_argument('--save-to-disk', dest='save_to_disk',
                        default="false",
                        help='Save to disk option')

    return parser.parse_args()


if __name__ == "__main__":
    main()