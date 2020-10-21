#!/bin/bash

# Copyright (c) 2019 Intel Corporation.
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

RED='\033[0;31m'
YELLOW="\033[1;33m"
GREEN="\033[0;32m"
NC='\033[0m' # No Color

function log_warn() {
    echo -e "${YELLOW}WARN: $1 ${NC}"
}

function log_info() {
    echo -e "${GREEN}INFO: $1 ${NC}"
}

function log_error() {
    echo -e "${RED}ERROR: $1 ${NC}"
}

function log_fatal() {
    echo -e "${RED}FATAL: $1 ${NC}"
    exit -1
}

function check_error() {
    if [ $? -ne 0 ] ; then
        log_fatal "$1"
    fi
}

cjson_version="1.7.12"

# URLs
cjson_url="https://github.com/DaveGamble/cJSON/archive/v${cjson_version}.tar.gz"

if [ ! -d "deps" ] ; then
    mkdir deps
    check_error "Failed to create dependencies directory"
fi

cd deps
check_error "Failed to change to dependencies directory"

# Installing cJSON dependency
if [ ! -f "cjson.tar.gz" ] ; then
    log_info "Downloading cJSON source"
    wget $cjson_url -O cjson.tar.gz
    check_error "Failed to download cJSON source"
fi

cjson_dir="cJSON-${cjson_version}"

if [ ! -d "$cjson_dir" ] ; then
    log_info "Extracting cJSON"
    tar xf cjson.tar.gz
    check_error "Failed to extract cJSON"
fi

cd $cjson_dir
check_error "Failed to change to cJSON directory"

if [ ! -d "build" ] ; then
    mkdir build
    check_error "Failed to create build directory"
fi

cd build
check_error "Failed to change to build directory"

log_info "Configuring cJSON for compilation"
cmake ..
check_error "Failed to configure cJSON"

log_info "Compiling cJSON library"
make -j$(nproc --ignore=2)
check_error "Failed to compile cJSON library"

log_info "Installing cJSON library"
make install
check_error "Failed to install cJSON library"

# Installing grpc dependency
# Library versions
grpc_version="v1.29.0"

# URLs
grpc_url="https://github.com/grpc/grpc"

# Dir to install grpc to
# Here, the prefix dir is /usr/local 
# and grpc lib gets installed to $grpc_install_prefix/lib
grpc_install_prefix="/usr/local"

# Dirs
grpc_dir="grpc"

if [ ! -d "deps" ] ; then
    mkdir deps
    check_error "Failed to create dependencies directory"
fi

cd deps
check_error "Failed to change to deps directory"

if [ ! -d "grpc" ] ; then
    log_info "git clone of grpc"
    git clone --recurse-submodules -b $grpc_version $grpc_url
    check_error "Failed to git clone"
fi

cd grpc
check_error "Failed to cd $grpc_dir"

if [ ! -d "cmake/build" ] ; then
    mkdir -p cmake/build
    check_error "Failed to create cmake/build directory"
fi

cd cmake/build
check_error "Falied to change directory to cmake/build"

log_info "Configuring lib grpc for building"
cmake -DCMAKE_INSTALL_PREFIX=$grpc_install_prefix ../..

check_error "Failed to configure lib grpc"

log_info "Compiling grpc library"
make
check_error "Failed to compile grpc library"

log_info "Installing grpc library"
make install
check_error "Failed to install grpc library"

log_info "Done."
