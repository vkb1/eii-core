# Copyright (c) 2019 Intel Corporation.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.



import queue
import threading
import imagestore_client
import time
import json
import influxdbconnector_client
import common
import sys
import os
import logging

logger=logging.getLogger()
logger.setLevel(logging.os.environ['PY_LOG_LEVEL'].upper())
is_done = False
def main():
    global is_done
    with open('../config/eis_config.json', 'r') as f:
        eis_config = json.load(f)

    with open('../config/query_config.json', 'r') as f:
        query_config = json.load(f)

    img_handle_queue = queue.Queue(maxsize=10000)

    #This thread will retriueve the image from imagestore service
    #and will store into frames directory
    img_rt_thread = threading.Thread(target = imagestore_client.retrieve_image_frames, args=(eis_config,query_config,img_handle_queue))
    img_rt_thread.daemon = True
    img_rt_thread.start()
    
    retrieve_measurement_data(eis_config,query_config,img_handle_queue)
    
    while True:
        if img_handle_queue.empty() and is_done:
            logger.info("Exiting...")
            sys.exit(0)
        else:
            time.sleep(5)

def retrieve_measurement_data(eis_config,query_config,img_handle_queue):
    global is_done
    influxdbconnector_client.query_influxdb(eis_config,query_config,img_handle_queue)
    is_done = True
    
if __name__ == "__main__":
    main()