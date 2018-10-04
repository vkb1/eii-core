/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

Explicit permissions are required to publish, distribute, sublicense, and/or sell copies of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE. 
*/

package streammanger

import (
	"fmt"
	"net"
	"os"
	"strings"
	"regexp"

	influxDBHelper "ElephantTrunkArch/Util"

	databus "ElephantTrunkArch/DataBusAbstraction/go"

	"github.com/golang/glog"
)

const (
	subscriptionName = "StreamManagerSubscription"
)

var opcuaContext = map[string]string{
	"direction": "PUB",
	"name":      "streammanager",
	"endpoint":  "opcua://%s:%s/elephanttrunk/",
}

// StrmMgr type
type StrmMgr struct {
	ServerHost        string //Server address
	ServerPort        string //Server Port
	InfluxDBHost      string
	InfluxDBPort      string
	InfluxDBName      string                     // Database Name
	MeasurementPolicy map[string]bool            //This DS to map policy name with action
	MsrmtTopicMap     map[string]OutStreamConfig //A map of mesurement to topic
	OpcuaPort         string
}

// OutStreamConfig type
type OutStreamConfig struct {
	RetenPolicy string // retention policy
	Measurement string // Measurement name
	Topic       string // Topic on databus
	StreamFlow  bool   // Whether to start the streaming or stop
	MsgBusType  string // Message bus type
}

func chkErr(err error) {
	//TODO: Will remove it with proper logging mechanism
	if err != nil {
		glog.Error("Stream Manager Error: ", err)
		// TODO: Yet to decide whether we can exit here or
		os.Exit(-1)
	}
}

func databPublish(dbus databus.DataBus, topic map[string]string, data string) error {

	err := dbus.Publish(topic, data)
	if err != nil {
		glog.Errorf("Publish Error: %v", err)
		return err
	}
	glog.Infof("Published [%s] : '%s'\n", topic, data)
	return nil

}

//This JsonConverter handels only the current line protocol which is specific.
//ToDo: Generic JsonConverter. To Handle all type of stream coming to StreamManager (Example: To handle Tags).

func JsonConverter(data string) string{
	var final_jbuf string
	final_jbuf="{\"Measurement\":" //Json string should start from { and 
				       //starting with a key "Measurement" for the value "classifier_results"

	jbuf := strings.Split(data," ") //Split the data based on single white-space

	final_jbuf+="\""+jbuf[0]+"\""+"," //Concatenating Measurement and classifier_results
					  // and comma seperated for the next key-value pair.

	for i:=2;i<=len(jbuf)-1;i++ {	// since number of white-spaces are there within the values 
		jbuf[1]+=" "+jbuf[i]    // of key-value pair, concatenating all the split string 
	}				//to one string to handle the further steps.

	influxTS:=",influx_ts="+jbuf[len(jbuf)-1]
	jbuf[1]=strings.Replace(jbuf[1],jbuf[len(jbuf)-1],influxTS,-1) //adding a key called influx_ts
								       //for the value of influx timestamp
	final_jbuf=final_jbuf+" "+jbuf[1]+"}"
	key_value_buf:=strings.Split(jbuf[1],"=")

	quoted_key:="\""+key_value_buf[0]+"\""
	final_jbuf=strings.Replace(final_jbuf,key_value_buf[0],quoted_key,-1) //wrapping "ImgHandle" key within quotes

	for j:=1;j<len(key_value_buf)-1;j++ {			//wrapping other keys with in the quotes
		key_buf:=strings.Split(key_value_buf[j],",")
		quoted_key2:="\""+key_buf[len(key_buf)-1]+"\""
		final_jbuf=strings.Replace(final_jbuf,key_buf[len(key_buf)-1],quoted_key2,-1)
	}

	final_jbuf=strings.Replace(final_jbuf,"=",":",-1)    //As the value of idx is in not in the Json
	regex:=regexp.MustCompile(`([0-9]+i)`)		     // acceptable format, hence making it string type
	final_jbuf=regex.ReplaceAllString(final_jbuf,`"$1"`) // by adding double quotes around it. 
	glog.Infof("Final Json String is = %s\n", final_jbuf)
	return final_jbuf
}

func (pStrmMgr *StrmMgr) handlePointData(dbus databus.DataBus, addr *net.UDPAddr, buf string, n int) error {
	var err error
	var Json_buf string
	err = nil
	glog.Infof("Received Data from address %s", addr)
	point := strings.Split(buf, " ")
	// It can have tag attached to it, let's split them too.
	msrTagsLst := strings.Split(point[0], ",")

	// Only allowing the measurements that are known to StreamManager
	for key, val := range pStrmMgr.MsrmtTopicMap {
		if key == msrTagsLst[0] {
			glog.Infof("Data Rcvd: %s\n", buf[0:n])
			glog.Infof("measurement and tag list: %s\n", msrTagsLst)
			// Publish only if a proper databus context available
			if dbus != nil {
				topic := map[string]string{
					"name": val.Topic,
					"type": "string",
				}

				Json_buf=JsonConverter(buf)

				databPublish(dbus, topic, Json_buf)
			}
		}
	}
	return err
}

// Init func subscribes to InfluxDB and starts up the udp server
func (pStrmMgr *StrmMgr) Init() error {

	//Setup the subscription for the DB
	// We have one DB only to be used by DA. Hence adding subscription
	// only during inititialization.

	client, err := influxDBHelper.CreateHTTPClient(pStrmMgr.InfluxDBHost,
		pStrmMgr.InfluxDBPort, "", "")

	if err != nil {
		glog.Errorf("Error creating InfluxDB client: %v", err)
	}

	defer client.Close()

	response, err := influxDBHelper.CreateSubscription(client, subscriptionName,
		pStrmMgr.InfluxDBName, pStrmMgr.ServerHost, pStrmMgr.ServerPort)

	if err == nil && response.Error() == nil {
		glog.Infoln("Successfully created subscription")
		go startServer(pStrmMgr)
	} else if response.Error() != nil {
		glog.Errorf("Response error: %v while creating subscription", response.Error())
		const str = "already exists"

		// TODO: we need to handle this situation in a more better way in
		// future in cases when DataAgent dies abruptly, system reboots etc.,
		if strings.Contains(response.Error().Error(), str) {
			glog.Infoln("subscription already exists, let's start the UDP" +
				" server anyways..")
			go startServer(pStrmMgr)
		}
	}

	return err
}

func startServer(pStrmMgr *StrmMgr) {

	var dstAddr string

	if pStrmMgr.ServerHost != "" {
		dstAddr = pStrmMgr.ServerHost + ":" + pStrmMgr.ServerPort
	} else {
		dstAddr = ":" + pStrmMgr.ServerPort
	}

	// Initialize the UDP server which listen on a predefined port
	servAddr, err := net.ResolveUDPAddr("udp", dstAddr)
	chkErr(err)

	servConnFd, err := net.ListenUDP("udp", servAddr)
	chkErr(err)
	glog.Infof("UDP server started gracefully at addr %s", dstAddr)
	defer servConnFd.Close()

	buf := make([]byte, 1024*1024)

	// TODO: check if this is the right place to put this
	opcuaDatab, err := databus.NewDataBus()
	if err != nil {
		glog.Errorf("New DataBus Instance creation Error: %v", err)
		os.Exit(1)
	}

	// This change is required to tie the opcua address to localhost or container's address
	hostname, err := os.Hostname()
	if err != nil {
		glog.Errorf("Failed to fetch the hostname of the node: %v", err)
	}
	if hostname != "ia_data_agent" {
		hostname = "localhost"
	}

	opcuaContext["endpoint"] = fmt.Sprintf(opcuaContext["endpoint"], hostname, pStrmMgr.OpcuaPort)

	err = opcuaDatab.ContextCreate(opcuaContext)
	if err != nil {
		glog.Errorf("DataBus-OPCUA context creation Error: %v", err)
		os.Exit(1)
	}
	defer opcuaDatab.ContextDestroy()
	glog.Infof("Databus-OPCUA context created")

	//TODO: break from for loop based on user signal,
	// Have a way to break from infinite loop
	//Dummy publish for OPCUA
	for _, val := range pStrmMgr.MsrmtTopicMap {
		if val.MsgBusType == "OPCUA" {
			dummyMsg := ""
			topic := map[string]string{"name": val.Topic, "type": "string"}
			err := opcuaDatab.Publish(topic, dummyMsg)
			if err != nil {
				glog.Errorf("Publish Error: %v", err)
			}
		}
	}
	for {
		//TODO: Should we break. Need channels here
		n, addr, err := servConnFd.ReadFromUDP(buf)
		if err != nil {
			glog.Errorln("ErrorSocketRecv: ", err)
			glog.Errorln("Stream manager: error in reading from socket")
		}

		// TODO: Need to use channels here for preventing the overflow of socket buffer
		pStrmMgr.handlePointData(opcuaDatab, addr, string(buf[:n]), n)
	}
}

// SetupOutStream func setups up out stream
func (pStrmMgr *StrmMgr) SetupOutStream(config *OutStreamConfig) error {

	var err error

	//populate the measurement->topic map
	pStrmMgr.MsrmtTopicMap[config.Measurement] = *config

	return err
}
