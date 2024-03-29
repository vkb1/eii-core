# Copyright (c) 2020 Intel Corporation.
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
"""EII Message Bus Server wrapper object
"""

import json

from .libeiiconfigmanager cimport *
from libc.stdlib cimport malloc
from libc.stdlib cimport free
from .util cimport Util


cdef class Server:
    """EII Message Bus Server object
    """

    def __init__(self):
        """Constructor
        """
        pass

    def __cinit__(self, *args, **kwargs):
        """Cython base constructor
        """
        self.app_cfg = NULL
        self.server_cfg = NULL

    @staticmethod
    cdef create(app_cfg_t* app_cfg, server_cfg_t* server_cfg):
        """Helper method for initializing the client object.

        :param app_cfg: Applications config struct
        :type: struct
        :param server_cfg: Server config struct
        :type: struct
        :return: Server class object
        :rtype: obj
        """
        s = Server()
        s.app_cfg = app_cfg
        s.server_cfg = server_cfg
        return s

    def __dealloc__(self):
        """Cython destructor
        """
        self.destroy()

    def destroy(self):
        """Destroy the server.
        """
        if self.server_cfg != NULL:
            server_cfg_config_destroy(self.server_cfg)

    def get_msgbus_config(self):
        """Constructs message bus config for Server

        :return: Messagebus config
        :rtype: dict
        """
        cdef char* config
        cdef config_t* msgbus_config
        try:
            msgbus_config = self.server_cfg.cfgmgr_get_msgbus_config_server(self.app_cfg.base_cfg, self.server_cfg)
            if msgbus_config is NULL:
                raise Exception("[Server] Getting msgbus config from base c layer failed")
        
            config = configt_to_char(msgbus_config)
            if config is NULL:
                raise Exception("[Server] config failed to get converted to char")

            config_str = config.decode('utf-8')
            free(config)
            config_destroy(msgbus_config)
            return json.loads(config_str)
        except Exception as ex:
            raise ex

    def get_interface_value(self, key):
        """To get particular interface value from Server interface config

        :param key: Key on which interface value will be extracted
        :type: string
        :return: Interface value
        :rtype: string
        """
        cdef config_value_t* value
        cdef char* config
        try:
            interface_value = None
            value = self.server_cfg.cfgmgr_get_interface_value_server(self.server_cfg, key.encode('utf-8'))
            if value is NULL:
                raise Exception("[Server] Getting interface value from base c layer failed")

            interface_value = Util.get_cvt_data(value)
            if interface_value is None:
                config_value_destroy(value)
                raise Exception("[Server] Getting cvt data failed")
        
            config_value_destroy(value)
            return interface_value
        except Exception as ex:
            raise ex

    def get_endpoint(self):
        """To get endpoint for particular server from its interface config

        :return: Endpoint config
        :rtype: string
        """
        cdef config_value_t* ep
        cdef char* c_endpoint
        try:
            ep = self.server_cfg.cfgmgr_get_endpoint_server(self.server_cfg)
            if ep is NULL:
                raise Exception("[Server] Getting end point from base c layer failed")

            if(ep.type == CVT_OBJECT):
                config = cvt_to_char(ep);
                if config is NULL:
                    config_value_destroy(ep)
                    raise Exception("[Server] Config cvt to char conversion failed")

                config_str = config.decode('utf-8')
                endpoint = json.loads(config_str)
            elif(ep.type == CVT_STRING):
                c_endpoint = ep.body.string
                if c_endpoint is NULL:
                    raise Exception("[Server] Endpoint getting string value failed")
                endpoint = c_endpoint.decode('utf-8')
            else:
                endpoint = None
                config_value_destroy(ep)
                raise TypeError("[Server] Type mismatch: EndPoint should be string or dict type")
        
            config_value_destroy(ep)
            return endpoint
        except TypeError as type_ex:
            raise type_ex
        except Exception as ex:
            raise ex

    def get_allowed_clients(self):
        """To get the names of the clients allowed to connect to server
        
        :return: List of clients
        :rtype: List
        """
        # Calling the base C cfgmgr_get_allowed_clients_server() API
        clients_list = []
        cdef config_value_t* clients
        cdef config_value_t* client_value
        cdef char* c_client_value
        try:
            clients = self.server_cfg.cfgmgr_get_allowed_clients_server(self.server_cfg)
            if clients is NULL:
                raise Exception("[Server] Getting alowed clients from base c layer failed")

            for i in range(config_value_array_len(clients)):
                client_value = config_value_array_get(clients, i)
                if client_value is NULL:
                    config_value_destroy(clients)
                    raise Exception("[Server] Getting array value from config for allowed clients failed.")

                c_client_value = client_value.body.string
                if c_client_value is NULL:
                    config_value_destroy(clients)
                    raise Exception("[Server] String value of client is NULL")

                client_val = c_client_value.decode('utf-8')
                clients_list.append(client_val)

                config_value_destroy(client_value)
            config_value_destroy(clients)
            return clients_list
        except Exception as ex:
            raise ex