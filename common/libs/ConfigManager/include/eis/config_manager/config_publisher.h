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

#ifndef _EIS_CH_CONFIGPUBLISHER_H
#define _EIS_CH_CONFIGPUBLISHER_H

#include <string.h>
#include <cjson/cJSON.h>
#include <functional> 
#include <iostream>
#include <vector>
#include <safe_lib.h>
#include <eis/utils/logger.h>
#include "eis/utils/json_config.h"
#include "db_client.h"
#include "eis/config_manager/config_app.h"

namespace eis {
    namespace config_manager {

        class PublisherCfg : public AppCfg {
            public:

                // App config
                config_t* config;

                // Publisher config
                config_value_t* publisher_cfg;

                /**
                * PublisherCfg Constructor
                */
                PublisherCfg(config_value_t* pub_config);

                /**
                 * Overridden base class method to fetch msgbus publisher configuration
                 * for application to communicate over EIS message bus
                 * @return config_t* - JSON msg bus server config of type config_t
                 */
                config_t* getMsgBusConfig() override;

                /**
                 * getEndpoint for application to fetch Endpoint associated with message bus config
                 * @return std::string - Endpoint of publisher config of type std::string
                 */
                std::string getEndpoint() override;

                /**
                 * getTopics for application to fetch the topics associated with message bus config
                 * @return vector<string> - Topics of publisher config
                 */
                std::vector<std::string> getTopics() override;

                /**
                 * setTopics for application to set topics associated with message bus config
                 * @return bool - Boolean whether topics were set
                 */
                bool setTopics(std::vector<std::string>) override;

                /**
                 * getAllowedClients for application to list of allowed clients associated with message bus config
                 * @return vector<string> - Allowed client of publisher config
                 */
                std::vector<std::string> getAllowedClients() override;

                // Destructor
                ~PublisherCfg();

        };
    }
}
#endif