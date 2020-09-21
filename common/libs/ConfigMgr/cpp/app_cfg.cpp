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
 * @brief AppCfg Implementation
 * Holds the implementaion of APIs supported by AppCfg class
 */

#include "eis/config_manager/app_cfg.hpp"
#include "eis/config_manager/cfg_mgr.h"

using namespace eis::config_manager;
using namespace std;


AppCfg::AppCfg(base_cfg_t* base_cfg) {
    m_base_cfg = base_cfg;
}

config_t* AppCfg::getConfig() {
    m_app_config = get_app_config(m_base_cfg);
    if (m_app_config == NULL) {
        LOG_ERROR_0("App Config not set");
        return NULL;
    }
    return m_app_config;
}

config_t* AppCfg::getInterface() {
    m_app_interface = get_app_interface(m_base_cfg);
    if (m_app_interface == NULL) {
        LOG_ERROR_0("App Interface not set");
        return NULL;
    }
    return m_app_interface;
}

config_value_t* AppCfg::getConfigValue(char* key) {
    config_value_t* value = get_app_config_value(m_base_cfg, key);
    if (value == NULL) {
        LOG_ERROR_0("Unable to fetch config value");
        return NULL;
    }
    return value;
}

config_value_t* AppCfg::getInterfaceValue(char* key) {
    config_value_t* value = get_app_interface_value(m_base_cfg, key);
    if (value == NULL) {
        LOG_ERROR_0("Unable to fetch interface value");
        return NULL;
    }
    return value;
}

// This virtual method is implemented
// by sub class objects
config_t* AppCfg::getMsgBusConfig() {

}

// This virtual method is implemented
// by sub class objects
std::string AppCfg::getEndpoint() {

}

// This virtual method is implemented
// by sub class objects
std::vector<std::string> AppCfg::getTopics() {

}

// This virtual method is implemented
// by sub class objects
std::vector<std::string> AppCfg::getAllowedClients() {

}

// This virtual method is implemented
// by sub class objects
bool AppCfg::setTopics(std::vector<std::string> topics_list) {

}

// tokenizer function to split string based on delimiter
vector<string> AppCfg::tokenizer(const char* str, const char* delim) {

    std::string line(str);
    std::stringstream str1(line);

    // Vector of string to save tokens
    vector<std::string> tokens;

    std::string temp;

    // Tokenizing w.r.t. delimiter
    while (getline(str1, temp, ':')) {
        tokens.push_back(temp);
    }

    return tokens;
}


AppCfg::~AppCfg() {
    if (m_app_config) {
        config_destroy(m_app_config);
    }
    if (m_app_interface) {
        config_destroy(m_app_interface);
    }
    if (m_app_data_store) {
        config_destroy(m_app_data_store);
    }
    if (m_base_cfg) {
        base_cfg_config_destroy(m_base_cfg);
    }
    LOG_INFO_0("ConfigMgr destructor");
}