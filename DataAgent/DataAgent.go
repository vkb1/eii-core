/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

Explicit permissions are required to publish, distribute, sublicense, and/or sell copies of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"flag"
	"os"

	config "ElephantTrunkArch/DataAgent/config"
	server "ElephantTrunkArch/DataAgent/da_grpc/server"
	internalserver "ElephantTrunkArch/DataAgent/da_grpc/server/server_internal"
	stm "ElephantTrunkArch/StreamManager"
	util "ElephantTrunkArch/Util"

	"github.com/golang/glog"
)

var strmMgrUDPServHost = "ia_data_agent"

const (
	strmMgrUDPServPort = "61971"
)

// DaCfg - stores parsed DataAgent config
var DaCfg config.DAConfig

func main() {

	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "config file path")

	flag.Parse()

	flag.Lookup("alsologtostderr").Value.Set("true")

	defer glog.Flush()
	if len(os.Args) < 2 {
		glog.Errorf("Usage: go run DataAgent/DataAgent.go " +
			"-config=<config_file_path> [-log_dir=<glog_dir_path>]")
		os.Exit(-1)
	}

	glog.Infof("**************STARTING DA**************")

	glog.Infof("Parsing the config file: %s", cfgPath)

	// Parse the DA config file
	err := DaCfg.ParseConfig(cfgPath)

	if err != nil {
		glog.Errorf("Error: %s while parsing config file: %s", err, cfgPath)
		os.Exit(-1)
	}

	influxCfg := DaCfg.InfluxDB
	clientadmin, err := util.CreateHTTPClient(influxCfg.Host,
		influxCfg.Port, "", "")

	if err != nil {
		glog.Errorf("Error creating InfluxDB client: %v", err)
		os.Exit(-1)
	}
	createadmin, err := util.CreateAdminUser(clientadmin, influxCfg.UserName, influxCfg.Password, influxCfg.DBName)

	if err == nil && createadmin.Error() == nil {
		glog.Infof("Successfully created admin user: %s", influxCfg.UserName)
	} else {
		if createadmin.Error() != nil {
			glog.Errorf("Error code: %v, Error Response: %s while creating "+
				"admin user: %s", err, createadmin.Error(), influxCfg.UserName)
		} else {
			glog.Errorf("Error code: %v while creating "+"admin user: %s", err, influxCfg.UserName)
		}
		os.Exit(-1)
	}
	clientadmin.Close()

	// Create InfluxDB database
	glog.Infof("Creating InfluxDB database: %s", influxCfg.DBName)
	client, err := util.CreateHTTPClient(influxCfg.Host,
		influxCfg.Port, influxCfg.UserName, influxCfg.Password)

	if err != nil {
		glog.Errorf("Error creating InfluxDB client: %v", err)
		os.Exit(-1)
	}

	response, err := util.CreateDatabase(client, influxCfg.DBName, influxCfg.Retention)
	if err != nil {
		glog.Infof("Cannot create database: %s", response.Error())
	}

	if err == nil && response.Error() == nil {
		glog.Infof("Successfully created database: %s", influxCfg.DBName)
	} else {
		if response.Error() != nil {
			glog.Errorf("Error code: %v, Error Response: %s while creating "+
				"database: %s", err, response.Error(), influxCfg.DBName)
		} else {
			glog.Errorf("Error code: %v while creating "+"database: %s", err, influxCfg.DBName)
		}
		os.Exit(-1)
	}
	defer client.Close()

	// Init StreamManager
	glog.Infof("**************STARTING STREAM MANAGER**************")

	var pStreamManager = new(stm.StrmMgr)

	// This change is required to tie the opcua address to localhost or container's address
	hostname, err := os.Hostname()
	if err != nil {
		glog.Errorf("Failed to fetch the hostname of the node: %v", err)
	}
	if strmMgrUDPServHost != hostname {
		strmMgrUDPServHost = "localhost"
	}

	pStreamManager.ServerHost = strmMgrUDPServHost
	pStreamManager.ServerPort = strmMgrUDPServPort
	pStreamManager.InfluxDBHost = DaCfg.InfluxDB.Host
	pStreamManager.InfluxDBPort = DaCfg.InfluxDB.Port
	pStreamManager.InfluxDBName = DaCfg.InfluxDB.DBName
	pStreamManager.InfluxDBUserName = DaCfg.InfluxDB.UserName
	pStreamManager.InfluxDBPassword = DaCfg.InfluxDB.Password
	pStreamManager.MsrmtTopicMap = make(map[string]stm.OutStreamConfig)
	pStreamManager.MeasurementPolicy = make(map[string]bool)
	pStreamManager.OpcuaPort = DaCfg.Opcua.Port

	glog.Infof("Going to start UDP server for influx subscription")
	err = pStreamManager.Init()
	if err != nil {
		glog.Errorf("Failed to initialize StreamManager : %v", err)
		os.Exit(-1)
	}

	var config = new(stm.OutStreamConfig)

	// Fetch the streams from the DA config file
	for key, val := range DaCfg.OutStreams {
		// TODO: Just using the 'key' as both measurement and topic for now.
		// This will change later
		config.Measurement = key
		config.Topic = key
		config.MsgBusType = val.DatabusFormat

		err = pStreamManager.SetupOutStream(config)
		if err != nil {
			glog.Errorf("Stream Manager, Error while setting up out stream: ", err)
		}
	}

	glog.Infof("**************STARTING GRPC SERVER**************")
	// Start GRPC server for GetConfig (Internal), GetConfig and GetQuery
	// external interfaces
	done := make(chan bool)
	glog.Infof("**************STARTING GRPC Internal SERVER**************")
	go internalserver.StartGrpcServer(DaCfg)
	glog.Infof("**************STARTING GRPC External SERVER**************")
	go server.StartGrpcServer(DaCfg)
	glog.Infof("**************Started GRPC servers**************")

	// Currently running this channel to keep the goroutine running
	// for StreamManager
	<-done
	glog.Infof("**************Exiting**************")
}
