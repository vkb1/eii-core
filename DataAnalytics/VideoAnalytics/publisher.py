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

import zmq
import threading
import logging
import os
import json
from concurrent.futures import ThreadPoolExecutor

class Publisher:

    def __init__(self, classifier_output_queue):
        """Constructor

        Parameters
        -----------
        classifier_output_queue : Queue
            Input queue for publisher (has [topic, classifier_metadata, 
            keyframe] data entries)
        """
        self.log = logging.getLogger(__name__)
        self.classifier_output_queue = classifier_output_queue
        self.stop_ev = threading.Event()

    def start(self):
        """
        Starts the publisher thread(s)
        """
        context = zmq.Context()
        socket = context.socket(zmq.PUB)
        topics = os.environ['PubTopics'].split(",")
        self.sockets = []
        self.publisher_threadpool = ThreadPoolExecutor(max_workers=len(topics))
        for topic in topics:
            topic_cfg = os.environ["{}_cfg".format(topic)].split(",")
            self.publisher_threadpool.submit(self.publish, socket, topic,
                                              topic_cfg)

    def publish(self, socket, topic, topic_cfg):
        """
        Send the data to the publish topic
        Parameters:
        ----------
        socket: ZMQ socket
            socket instance
        topic: str
            topic name
        topic_cfg: str
            topic config
        """
        thread_id = threading.get_ident()
        self.log.info("Publisher thread ID: {} started" +  \
                      " with topic: {} and topic_cfg: {}...".format(thread_id,
                      topic, topic_cfg))

        mode = topic_cfg[0].lower()
        try:
            if "tcp" in mode:
                socket.bind("tcp://{}".format(topic_cfg[1]))
            elif "ipc" in mode:
                socket.bind("ipc://{}".format(topic_cfg[1]))
            self.sockets.append(socket)
        except Exception as ex:
            self.log.exception(ex)

        self.log.info("Publishing to topic: {}...".format(topic))
        
        while not self.stop_ev.is_set():
            result = self.classifier_output_queue.get()
            try:
                result[0] = topic.encode()
                result[1] = json.dumps(result[1]).encode()                
                socket.send_multipart(result, copy=False)
                self.log.debug("Published data : {}...".format(result[0:2]))
            except Exception as ex:
                self.log.exception('Error while publishing data: {}'.format(ex))
        self.log.info("Publisher thread ID: {} stopped" +  \
                      " with topic: {} and topic_cfg: {}...".format(thread_id,
                      topic, topic_cfg))

    def stop(self):
        """
        Stops the Subscriber thread
        """
        for socket in self.sockets:
            socket.close()
            if socket._closed == "False":
                self.log.error("Unable to close socket connection")
        self.stop_ev.set()
        self.publisher_threadpool.shutdown(wait=False)
