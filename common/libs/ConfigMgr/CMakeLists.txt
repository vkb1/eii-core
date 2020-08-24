# Copyright (c) 2019 Intel Corporation.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to
# deal in the Software without restriction, including without limitation the
# rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
# sell copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
# FROM,OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
# IN THE SOFTWARE.

cmake_minimum_required(VERSION 3.11)
project(eis-configmanager VERSION 1.0.0 LANGUAGES C CXX)

# Define CMake options

option(WITH_EXAMPLES "Compile with examples" ON)

# Globals
set(EIS_COMMON_CMAKE "${CMAKE_CURRENT_SOURCE_DIR}/../../cmake")

# Set C standard to C11
set(CMAKE_CXX_STANDARD 11)

# Set CFLAGS
set(CMAKE_C_FLAGS "-fPIE -fPIC -O2 -Wall -pedantic -fstack-protector-strong -D_FORTIFY_SOURCE=2")

# Set LDFLAGS
set(CMAKE_SHARED_LINKER_FLAGS "-z noexecstack -z relro -z now")

# Add cmake module path
set(CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/cmake" "${EIS_COMMON_CMAKE}")

# Finding libraries
find_package(cJSON REQUIRED)
find_package(EISUtils REQUIRED)
find_package(EISMsgEnv REQUIRED)
find_package(IntelSafeString REQUIRED)

# Include header directories
include_directories(
    include/
    db_client/include/
    db_client/etcd_client/include/
    ${EISMsgEnv_INCLUDE}
    ${EISUtils_INCLUDE}
    ${IntelSafeString_INCLUDE})

# Get all source files
file(GLOB SOURCES "src/*.c" "src/*.cpp")
set_source_files_properties(${SOURCES} PROPERTIES LANGUAGE C)

add_library(neweisconfigmgr SHARED ${SOURCES})

target_link_libraries(neweisconfigmgr
    PUBLIC
        pthread
        ${EISMsgEnv_LIBRARIES}
        ${EISUtils_LIBRARIES}
        eisdbclient
    PRIVATE
        ${IntelSafeString_LIBRARIES})

# If compile in debug mode, set DEBUG flag for C code
if("${CMAKE_BUILD_TYPE}" STREQUAL "Debug")
    target_compile_definitions(neweisconfigmgr PRIVATE DEBUG=1)
endif()

if(WITH_EXAMPLES)
    add_subdirectory(examples/)
endif()

##
## Configure pkg-config file to be installed for the EIS ConfigMgr
##
set(PKG_CONFIG_IN  "${CMAKE_CURRENT_SOURCE_DIR}/cmake/libneweisconfigmgr.pc.in")
set(PKG_CONFIG_OUT "${CMAKE_CURRENT_BINARY_DIR}/libneweisconfigmgr.pc")
set(DEST_DIR       "${CMAKE_INSTALL_PREFIX}")
set(PRIVATE_LIBS   "-lzmq -lcjson")

configure_file(${PKG_CONFIG_IN} ${PKG_CONFIG_OUT} @ONLY)

##
## Add CMake configuration for installing the library including files for other
## projects finding the library using CMake
##

include(GNUInstallDirs)
set(INSTALL_CONFIGDIR ${CMAKE_INSTALL_LIBDIR}/cmake/EISConfigMgrNew)

install(TARGETS neweisconfigmgr
    EXPORT neweisconfigmgr-targets
    LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
    ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR})

set_target_properties(neweisconfigmgr PROPERTIES EXPORT_NAME EISConfigMgrNew)
install(DIRECTORY include/ DESTINATION ${CMAKE_INSTALL_INCLUDEDIR})

# Install pkg-config libneweisconfigmgr.pc file
install(
    FILES
        ${PKG_CONFIG_OUT}
    DESTINATION
        ${CMAKE_INSTALL_LIBDIR}/pkgconfig
)

# Export targets to a script
install(EXPORT neweisconfigmgr-targets
    FILE
        EISConfigMgrNewTargets.cmake
    DESTINATION
        ${INSTALL_CONFIGDIR}
)

# Create a ConfigVersion.cmake file
include(CMakePackageConfigHelpers)
write_basic_package_version_file(
    ${CMAKE_CURRENT_BINARY_DIR}/EISConfigMgrNewConfigVersion.cmake
    VERSION ${PROJECT_VERSION}
    COMPATIBILITY AnyNewerVersion
)

configure_package_config_file(
    ${CMAKE_CURRENT_LIST_DIR}/cmake/EISConfigMgrNewConfig.cmake.in
    ${CMAKE_CURRENT_BINARY_DIR}/EISConfigMgrNewConfig.cmake
    INSTALL_DESTINATION ${INSTALL_CONFIGDIR}
)

# Install the config, configversion and custom find modules
install(FILES
    ${CMAKE_CURRENT_BINARY_DIR}/EISConfigMgrNewConfigVersion.cmake
    ${CMAKE_CURRENT_BINARY_DIR}/EISConfigMgrNewConfig.cmake
    DESTINATION ${INSTALL_CONFIGDIR}
)

export(EXPORT neweisconfigmgr-targets
       FILE ${CMAKE_CURRENT_BINARY_DIR}/EISConfigMgrNewTargets.cmake)

# Register package in user's package registry
export(PACKAGE EISConfigMgrNew)