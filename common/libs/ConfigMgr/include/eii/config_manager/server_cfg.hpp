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
// FROM,OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.


/**
 * @file
 * @brief ConfigMgr interface
 */

#ifndef _EIS_CH_SERVER_CFG_H
#define _EIS_CH_SERVER_CFG_H

#include <string.h>
#include <cjson/cJSON.h>
#include <iostream>
#include <safe_lib.h>
#include <eis/utils/logger.h>
#include "eis/utils/json_config.h"
#include "eis/config_manager/kv_store_plugin/kv_store_plugin.h"
#include "eis/config_manager/app_cfg.hpp"
#include "eis/config_manager/util_cfg.h"
#include "eis/config_manager/cfg_mgr.h"


namespace eis {
    namespace config_manager {

        class ServerCfg : public AppCfg {
            private:
                // server_cfg_t object
                server_cfg_t* m_serv_cfg;

                // app_cfg_t object
                app_cfg_t* m_app_cfg;
            public:
                /**
                * ServerCfg Constructor
                * @param server_config - The config associated with a server
                */
                explicit ServerCfg(server_cfg_t* serv_cfg, app_cfg_t* app_cfg);

                /**
                 * Overridden base class method to fetch msgbus server configuration
                 * for application to communicate over EIS message bus
                 * @return config_t* - JSON msg bus server config of type config_t
                 */
                config_t* getMsgBusConfig() override;

                /**
                 * Overridden base class method to fetch interface value
                 * for application to communicate over EIS message bus
                 * @param key - Key on which interface value is extracted.
                 * @return config_value_t* - config_value_t object
                 */
                config_value_t* getInterfaceValue(const char* key) override;

                /**
                 * getEndpoint for application to fetch Endpoint associated with message bus config
                 * @return std::string - Endpoint of server config of type std::string
                 */
                std::string getEndpoint() override;

                /**
                 * getAllowedClients for application to list of allowed clients associated with message bus config
                 * @return vector<string> - Allowed client of server config
                 */
                std::vector<std::string> getAllowedClients() override;

                /**
                * server_cfg_t getter to get private m_serv_cfg
                */
                server_cfg_t* getServCfg();

                /**
                * app_cfg_t getter to get private m_app_cfg
                */
                app_cfg_t* getAppCfg();

                /**
                * Destructor
                */
                ~ServerCfg();

        };
    }
}
#endif
