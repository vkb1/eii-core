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
 * @brief ConfigMgr AppCfg class
 */

#ifndef _EII_CH_APP_CFG_H
#define _EII_CH_APP_CFG_H

#include <string.h>
#include <cjson/cJSON.h>
#include <iostream>
#include <string>
#include <vector>
#include <bits/stdc++.h>
#include <safe_lib.h>
#include <eii/utils/logger.h>
#include "eii/utils/json_config.h"
#include "eii/config_manager/kv_store_plugin/kv_store_plugin.h"
#include "eii/config_manager/base_cfg.h"


namespace eii {
    namespace config_manager {

        /**
         * AppCfg class
         */  
        class AppCfg {
            private:

                // App's config
                config_t* m_app_config;

                // App's interface
                config_t* m_app_interface;

                // C base_cfg_t struct
                base_cfg_t* m_base_cfg;

                // App's data store
                config_t* m_app_data_store;

            protected:

                /**
                 * Helper base class function to split string based on delimiter
                 * @param str - string to be split
                 * @param delim - delimiter
                 * @return std::vector<std::string> - vector of split strings
                 */
                std::vector<std::string> tokenizer(const char* str,
                                                   const char* delim);

            public:

                /**
                * AppCfg Constructor
                * @param app_config - The config associated with a service
                * @param app_interface - The interface associated with a service
                * @param dev_mode - bool whether dev mode is set
                */
                explicit AppCfg(base_cfg_t* base_cfg);

                /**
                 * Gets app config
                 * @param key - Key for which value is needed
                 * @return config_value_t* - config_value_t object
                 */
                config_t* getConfig();

                /**
                 * Gets app interface
                 * @return config_value_t* - config_value_t object
                 */
                config_t* getInterface();

                /**
                 * Gets value from respective application's config
                 * @param key - Key for which value is needed
                 * @return config_value_t* - config_value_t object
                 */
                config_value_t* getConfigValue(char* key);

                /**
                 * Register a callback to watch on any given key
                 * @param key - key to watch
                 * @param watch_callback - callback object
                 * @param user_data - user data to be sent to callback
                 * @return bool - Boolean whether the callback was registered
                 */
                bool watch(char* key, callback_t watch_callback, void* user_data);

                /**
                 * Register a callback to watch on any given key prefix
                 * @param prefix - key prefix to watch
                 * @param watch_callback - callback object
                 * @param user_data - user data to be sent to callback
                 * @return bool - Boolean whether the callback was registered
                 */
                bool watchPrefix(char* prefix, callback_t watch_callback, void* user_data);

                /**
                 * Register a callback to watch on app config
                 * @param watch_callback - callback object
                 * @param user_data - user data to be sent to callback
                 * @return bool - Boolean whether the callback was registered
                 */
                bool watchConfig(callback_t watch_callback, void* user_data);

                /**
                 * Register a callback to watch on app interface
                 * @param watch_callback - callback object
                 * @param user_data - user data to be sent to callback
                 * @return bool - Boolean whether the callback was registered
                 */
                bool watchInterface(callback_t watch_callback, void* user_data);

                /**
                 * Get msgbus configuration for application to communicate over EII message bus
                 * @return config_t* - JSON msg bus server config of type config_t
                 *                   - On Failure, returns NULL
                 */ 
                virtual config_t* getMsgBusConfig();

                /**
                 * virtual getEndpoint function implemented by child classes to fetch Endpoint
                 * @return std::string - On Success, Endpoint of associated config of type std::string
                 *                     - On Failure, returns empty string
                 */
                virtual std::string getEndpoint();

                /**
                 * virtual function that gets value from interface
                 * @param key - Key for which value is needed
                 * @return config_value_t* - On Success, config_value_t object
                 *                         - On Failure, returns NULL
                 */
                virtual config_value_t* getInterfaceValue(const char* key);

                /**
                 * virtual getTopics function implemented by child classes to fetch topics
                 * @return std::string - On Success, Topics of associated config of type std::string
                 *                     - On Failure, returns empty vector
                 */
                virtual std::vector<std::string> getTopics();

                /**
                 * virtual setTopics function implemented by child classes to set topics
                 * @param topics_list - vector of strings containing topics
                 * @return bool - Boolean whether topics were set
                 *              - On Success, returns true
                 *              - On Failure, returns false
                 */
                virtual bool setTopics(std::vector<std::string> topics_list);

                /**
                 * virtual getAllowedClients function implemented by child classes to fetch topics
                 * @return std::string - On Success, Allowed client of associated config of type std::string
                 *                     - On Failure, returns empty vector
                 */
                virtual std::vector<std::string> getAllowedClients();

                /**
                * Destructor
                */
                virtual ~AppCfg();
        };
    }
}
#endif
