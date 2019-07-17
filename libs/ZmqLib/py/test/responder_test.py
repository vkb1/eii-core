import cv2
import os
import argparse
import json
import etcd3
from distutils.util import strtobool
from libs.ZmqLib.py.rep.responder import ZMQ_RESPONDER
from libs.ConfigManager.etcd.py.etcd_client import EtcdCli


def callback_a(request):
    print("Request from client {}".format(request))
    args = parse_args()
    save_to_disk = bool(strtobool(args.save_to_disk))
    etcd = etcd3.client(host="localhost", port=2379,
                        ca_cert="/run/secrets/etcd.ca.cert",
                        cert_key="/run/secrets/etcd.client.key",
                        cert_cert="/run/secrets/etcd.client.cert")

    # Initiate the OpenCV VideoCapture
    if not os.path.exists("./pcb_d2000.avi"):
        print("Video file does not exist...")
    cap = cv2.VideoCapture("./pcb_d2000.avi")
    if not cap.isOpened():
        print("Failed to open video file...")

    # Read frames from video and send meta-data and numpy frame
    ret, frame = cap.read()
    if ret:
        metaData = dict(
            dtype=str(frame.dtype),
            shape=frame.shape,
            name="{0}".format('Frame')
        )

        # Sample etcd tests
        result = json.dumps(metaData)
        etcd.put('metaData', result)
        etcd_metaData = etcd.get('metaData')
        etcd.delete('metaData')
        send_list = [etcd_metaData[0].decode("utf-8"),
                     frame]

        # Write sent image frame to disk
        if save_to_disk:
            outputFilePath = "./{0}".format(metaData['name'])
            with open(outputFilePath, "wb") as outfile:
                outfile.write(frame)
        cap.release()
        return send_list
    print("Failed to read video file...")


def main():
    """main method"""
    args = parse_args()
    etcd = etcd3.client(host="localhost", port=2379,
                        ca_cert="/run/secrets/etcd.ca.cert",
                        cert_key="/run/secrets/etcd.client.key",
                        cert_cert="/run/secrets/etcd.client.cert")

    config_list = os.environ["connection_config"].split(',')
    responder_address = ""
    if "tcp" in config_list[0].lower():
        responder_address = "tcp://"+config_list[1]
    elif "ipc" in config_list[0].lower():
        responder_address = "ipc:///"+config_list[1]

    with open(args.public_keys_dir+"/client.key", "rb") as client_key:
        client_key_encoded = client_key.read()
    with open(args.public_keys_dir+"/server.key", "rb") as server_key:
        server_key_encoded = server_key.read()

    etcd.put('server_key', server_key_encoded)
    etcd.put('client_key', client_key_encoded)

    conf = {"endpoint": "localhost:2379",
            "certFile": "/run/secrets/etcd.client.cert",
            "keyFile": "/run/secrets/etcd.client.key",
            "trustFile": "/run/secrets/etcd.ca.cert"}
    etcdCli = EtcdCli(conf)
    server_etcd_key = etcdCli.GetConfig('server_key')
    client_etcd_key = etcdCli.GetConfig('client_key')
    etcd_keys_directory = os.path.join(os.getcwd(), r'etcd_public_keys')
    if not os.path.exists(etcd_keys_directory):
        os.makedirs(etcd_keys_directory)
    with open(etcd_keys_directory+"/server.key", "wb") as server_key_file:
        server_key_file.write(server_etcd_key.encode('utf-8'))
    with open(etcd_keys_directory+"/client.key", "wb") as client_key_file:
        client_key_file.write(client_etcd_key.encode('utf-8'))

    # Initiate the responder
    responder = ZMQ_RESPONDER(responder_address,
                              public_keys_dir=args.public_keys_dir,
                              secret_keys_dir=args.private_keys_dir)

    while True:
        # Send the meta-data and numpy frame
        responder.send(callback=callback_a)

    # Close socket
    responder.close_socket()


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