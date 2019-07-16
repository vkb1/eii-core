from libs.ZmqLib.py.req.requester import ZMQ_REQUESTER
import time
from distutils.util import strtobool
import argparse
import os


def main():
    """ main method """
    args = parse_args()
    save_to_disk = bool(strtobool(args.save_to_disk))

    config_list = os.environ["connection_config"].split(',')
    request_address = ""
    if "tcp" in config_list[0].lower():
        request_address = "tcp://"+config_list[1]
    elif "ipc" in config_list[0].lower():
        request_address = "ipc:///"+config_list[1]

    # Initiate the requester
    requester = ZMQ_REQUESTER(request_address,
                              public_keys_dir=args.public_keys_dir,
                              secret_keys_dir=args.private_keys_dir)
    start_time = time.time()
    frame_count = 0
    while True:

        # Receive on a loop
        meta_data, frame = requester.receive("imgHandle")

        frame_count += 1
        if time.time() - start_time >= 1:
            print("FPS is : {}".format(frame_count))
            frame_count = 0
            start_time = time.time()

        # Write obtained image frame to disk
        if save_to_disk:
            outputFilePath = "./{0}{1}".format(meta_data['name'], "_out")
            with open(outputFilePath, "wb") as outfile:
                outfile.write(frame)

    # Close socket
    requester.close_socket()


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
