// Copyright (c) 2020 Intel Corporation.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

/**
 * @file
 * @brief Publisher config implmentaion
 */


#include "eis/config_manager/base_cfg.h"
#define TYPE "Type"
#define NAME "Name"
#define ENDPOINT "EndPoint"
#define PUBLISHERS "Publishers"
#define ZMQ_RECV_HWM "zmq_recv_hwm"
#define ALLOWED_CLIENTS "AllowedClients"
#define PUBLIC_KEYS "/Publickeys/"
#define PRIVATE_KEY "/private_key"

#ifndef _EIS_C_PUB_CFG_H
#define _EIS_C_PUB_CFG_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {

    config_t* (*cfgmgr_get_msgbus_config_pub)(base_cfg_t* base_cfg);

    config_value_t* (*cfgmgr_get_interface_value_pub)(base_cfg_t* base_cfg, const char* key);
    
    config_value_t* (*cfgmgr_get_endpoint_pub)(base_cfg_t* base_cfg);

    config_value_t* (*cfgmgr_get_topics_pub)(base_cfg_t* base_cfg);

    int (*cfgmgr_set_topics_pub)(char** topics_list, int len, base_cfg_t* base_cfg);

    config_value_t* (*cfgmgr_get_allowed_clients_pub)(base_cfg_t* base_cfg);

    config_value_t* pub_config;

} pub_cfg_t;

/**
 * pub_cfg_new function to creates a new pub_cfg_t object
 *  @return NULL for any errors occured or pub_cfg_t* on successful
 */
pub_cfg_t* pub_cfg_new();

/**
 * Destroy pub_cfg_t* object.
 *
 * @param pub_cfg_config - configuration to destroy
 */
void pub_cfg_config_destroy(pub_cfg_t *pub_cfg_config);

#ifdef __cplusplus
}
#endif

#endif