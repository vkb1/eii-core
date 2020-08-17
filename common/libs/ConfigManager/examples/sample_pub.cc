// Copyright (c) 2019 Intel Corporation.
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
 * @brief Message bus publisher example
 * @author Kevin Midkiff (kevin.midkiff@intel.com)
 */

#include <signal.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "eis/msgbus/msgbus.h"
#include "eis/utils/logger.h"
#include "eis/utils/json_config.h"
#include "eis/config_manager/config_mgr.h"

#define TOPIC "publish_test"

using namespace eis::config_manager;

// Globals for cleaning up nicely
publisher_ctx_t* g_pub_ctx = NULL;
msg_envelope_t* g_msg = NULL;
void* g_msgbus_ctx = NULL;

/**
 * Helper to initailize the message to be published
 */
void initialize_message() {
    // Creating message to be published
    msg_envelope_elem_body_t* integer = msgbus_msg_envelope_new_integer(42);
    msg_envelope_elem_body_t* fp = msgbus_msg_envelope_new_floating(55.5);
    g_msg = msgbus_msg_envelope_new(CT_JSON);
    msgbus_msg_envelope_put(g_msg, "hello", integer);
    msgbus_msg_envelope_put(g_msg, "world", fp);
}

/**
 * Function to print publisher usage
 */
void usage(const char* name) {
    fprintf(stderr, "usage: %s [-h|--help] <json-config> [topic]\n", name);
    fprintf(stderr, "\t-h|--help   - Show this help\n");
    fprintf(stderr, "\tjson-config - Path to JSON configuration file\n");
    fprintf(stderr, "\ttopic       - (Optional) Topic string "\
                    "(df: publish_test)\n");
}

/**
 * Signal handler
 */
void signal_handler(int signo) {
    LOG_INFO_0("Cleaning up");
    if(g_pub_ctx != NULL) {
        LOG_INFO_0("Freeing publisher");
        msgbus_publisher_destroy(g_msgbus_ctx, g_pub_ctx);
        g_pub_ctx = NULL;
    }
    if(g_msg != NULL) {
        LOG_INFO_0("Freeing message");
        msgbus_msg_envelope_destroy(g_msg);
        g_msg = NULL;
    }
    if(g_msgbus_ctx != NULL) {
        LOG_INFO_0("Freeing message bus context");
        msgbus_destroy(g_msgbus_ctx);
        g_msgbus_ctx = NULL;
    }
    LOG_INFO_0("Done.");
}

int main(int argc, char** argv) {

    // Set log level
    set_log_level(LOG_LVL_DEBUG);

    // Setting up signal handlers
    signal(SIGINT, signal_handler);
    signal(SIGTERM, signal_handler);

    setenv("DEV_MODE", "TRUE", 1);

    // Fetching Publisher config from
    // VideoIngestion interface
    setenv("AppName","VideoIngestion", 1);
    ConfigMgr* pub_ch = new ConfigMgr();

    PublisherCfg* pub_ctx = pub_ch->getPublisherByIndex(0);
    config_t* pub_config = pub_ctx->getMsgBusConfig();

    // Initializing Publisher using pub_config obtained
    // from new ConfigManager APIs
    g_msgbus_ctx = msgbus_initialize(pub_config);
    if(g_msgbus_ctx == NULL) {
        LOG_ERROR_0("Failed to initialize message bus");
        goto err;
    }

    msgbus_ret_t ret;
    ret = msgbus_publisher_new(g_msgbus_ctx, TOPIC, &g_pub_ctx);
    if(ret != MSG_SUCCESS) {
        LOG_ERROR("Failed to initialize publisher (errno: %d)", ret);
        goto err;
    }

    // Initialize message to be published
    initialize_message();

    LOG_INFO_0("Running...");
    while(g_pub_ctx != NULL) {
        LOG_INFO_0("Publishing message");
        ret = msgbus_publisher_publish(g_msgbus_ctx, g_pub_ctx, g_msg);
        if(ret != MSG_SUCCESS) {
            LOG_ERROR("Failed to publish message (errno: %d)", ret);
            goto err;
        }
        sleep(1);
    }

    return 0;

err:
    if(g_pub_ctx != NULL)
        msgbus_publisher_destroy(g_msgbus_ctx, g_pub_ctx);
    if(g_msgbus_ctx != NULL)
        msgbus_destroy(g_msgbus_ctx);
    return -1;
}
