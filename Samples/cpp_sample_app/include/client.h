#include <signal.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <atomic>

#include "eis/msgbus/msgbus.h"
#include "eis/utils/logger.h"
#include "eis/utils/json_config.h"
#include <eis/config_manager/env_config.h>

#ifndef _EIS_CPP_CLIENT_H
#define _EIS_CPP_CLIENT_H

class Client
{
    private:
        env_config_t* g_env_config_client = NULL;
        config_mgr_t* g_config_mgr = NULL;
        void* g_msgbus_ctx_client = NULL;
        recv_ctx_t* g_service_ctx = NULL;
        msg_envelope_serialized_part_t* parts = NULL;
        msg_envelope_t* msg = NULL;
        int num_parts = 0;
        std::atomic<bool> *loop;

    public:
        Client(std::atomic<bool> *loop);
        ~Client();
        bool init(char *service_name);
        static void* start(void *arg);
        int client_service();
        void clean_up();
};

#endif //_EIS_CPP_CLIENT_H

