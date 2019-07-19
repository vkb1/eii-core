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
// FROM,OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

/**
 * @brief Logging primitives
 * @author Kevin Midkiff <kevin.midkiff@intel.com>
 */

#ifndef _EIS_MESSAGE_BUS_LOG_H
#define _EIS_MESSAGE_BUS_LOG_H

#include <time.h>
#include <stdio.h>

#ifdef __cplusplus
extern "C" {

#endif

#define LOG(lvl, fmt, ...) { \
    time_t ltime; \
    ltime = time(NULL); \
    char* t = asctime(localtime(&ltime)); \
    t[24] = '\0'; \
    fprintf(stderr, "[%s] %5s:%s:%d: " fmt "\n", t, \
            lvl, __func__, __LINE__, ##__VA_ARGS__); \
}
#define LOG_0(lvl, msg) { \
    time_t ltime; \
    ltime = time(NULL); \
    char* t = asctime(localtime(&ltime)); \
    t[24] = '\0'; \
    fprintf(stderr, "[%s] %5s:%s:%d: " msg "\n", t, \
            lvl, __func__, __LINE__); \
}

#ifdef DEBUG
#define LOG_DEBUG(fmt, ...) LOG("DEBUG", fmt, ##__VA_ARGS__)
#define LOG_DEBUG_0(msg) LOG_0("DEBUG", msg)
#else
#define LOG_DEBUG(fmt, ...)
#define LOG_DEBUG_0(msg)
#endif

#define LOG_INFO(fmt, ...) LOG("INFO", fmt, ##__VA_ARGS__)
#define LOG_INFO_0(msg) LOG_0("INFO", msg)
#define LOG_WARN(fmt, ...) LOG("WARN", fmt, ##__VA_ARGS__)
#define LOG_WARN_0(msg) LOG_0("WARN", msg)
#define LOG_ERROR(fmt, ...) LOG("ERROR", fmt, ##__VA_ARGS__)
#define LOG_ERROR_0(msg) LOG_0("ERROR", msg)

#ifdef __cplusplus
}
#endif

#endif // _EIS_MESSAGE_BUS_LOG_H
