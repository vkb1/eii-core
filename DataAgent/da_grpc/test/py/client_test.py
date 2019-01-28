"""
Copyright (c) 2018 Intel Corporation.
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
"""

# Python grpc client implementation

import logging
import argparse
from ImageStore.client.py.client import GrpcImageStoreClient
import hashlib
import time
import sys
import os
from DataAgent.da_grpc.client.py.client \
    import GrpcClient

logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s : %(levelname)s : \
                    %(name)s : [%(filename)s] :' +
                    '%(funcName)s : in line : [%(lineno)d] : %(message)s')
log = logging.getLogger("GRPC_TEST")


CLIENT_CERT = "/etc/ssl/grpc_int_ssl_secrets/client_certificate.pem"
CLIENT_KEY = "/etc/ssl/grpc_int_ssl_secrets/client_key.pem"
CA_CERT = "/etc/ssl/grpc_int_ssl_secrets/ca_certificate.pem"

IM_CLIENT_KEY = "/etc/ssl/imagestore/imagestore_client_certificate.pem"
IM_CLIENT_CERT = "/etc/ssl/imagestore/imagestore_client_key.pem"


def parse_args():
    """Parse command line arguments
    """
    parser = argparse.ArgumentParser()
    parser.add_argument('--input_file', help='input image file')
    parser.add_argument('--output_file', help='output image file')

    return parser.parse_args()

if __name__ == '__main__':

    args = parse_args()

    # If executing this script from other m/c, provide
    # the right hostname/ip addr of the system running
    # DataAgent module of IEI
    client = GrpcClient(CLIENT_CERT, CLIENT_KEY, CA_CERT, hostname="localhost")

    # Testing GetBlob(imgHandle) gRPC call
    inputFile = args.input_file
    outputFile = args.output_file

    imgStore = GrpcImageStoreClient(IM_CLIENT_CERT, IM_CLIENT_KEY, CA_CERT)

    totalTime3 = 0.0
    iter1 = 20
    for index in range(iter1):
        log.info("Binary reading of file: %s", inputFile)
        inputBytes = None
        with open(inputFile, "rb") as f:
            inputBytes = f.read()

        log.info("len(inputBytes): %s", len(inputBytes))
        log.info("Storing the binary data read in ImageStore...")
        try:
            key = imgStore.Store(inputBytes, "inmemory")
        except Exception as ex:
            log.error("Error while doing imgStore.store operation...")
            log.error(ex)
            continue

        log.info("Image Handle obtainer after imgStore.store() operation: %s",
                 key)

        log.info("Calling GetBlob(%s) gRPC interface...", key)
        start = time.time()
        outputBytes = client.GetBlob(key)
        end = time.time()
        timeTaken = end - start
        totalTime3 += timeTaken

        log.info("len(outputBytes): %s", len(outputBytes))

        log.info("Writing the binary data received into a file: %s",
                 outputFile)
        with open(outputFile, "wb") as outfile:
            outfile.write(outputBytes)

        digests = []
        for filename in [inputFile, outputFile]:
            hasher = hashlib.md5()
            with open(filename, 'rb') as f:
                buf = f.read()
                hasher.update(buf)
                a = hasher.hexdigest()
                digests.append(a)
                log.info("Hash for filename: %s is %s", filename, a)

        if digests[0] == digests[1]:
            log.info("md5sum for the files match")
        else:
            log.info("md5sum for the files doesn't match")

        log.info("GetBlob call...index: %d, time: %f secs", index, timeTaken)

    log.info("Average time taken for GetBlob() %d calls: %f secs",
             iter1, totalTime3 / iter1)
