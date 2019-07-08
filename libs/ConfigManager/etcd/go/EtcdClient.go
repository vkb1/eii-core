/*
Copyright (c) 2019 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// Package etcdclient is a etc wrapper to read any keys and watch on a directory in etcd distributed key store
package etcdclient

import (
	util "IEdgeInsights/Util"
	"context"
	"errors"
	"strings"

	"github.com/golang/glog"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
)

// Config struct
type Config struct {
	Endpoint  []string
	CertFile  string
	KeyFile   string
	TrustFile string
}

// EtcdClient interface lets your Get/Set from Etcd
type EtcdClient interface {
	GetConfig(key string) (string, error)

	RegisterDirWatch(key string, onChangeCallback OnChangeCallback)

	RegisterKeyWatch(key string, onChangeCallback OnChangeCallback)
}

// OnChangeCallback is used for passing callbacks to
// RegisterDirWatch and RegisterKeyWatch go routines
type OnChangeCallback func(key string, newValue string)

// EtcdCli implements EtcdClient
type EtcdCli struct {
	etcd *clientv3.Client
}

// RegisterWatchOnKey go routine
func RegisterWatchOnKey(rch clientv3.WatchChan, onChangeCallback OnChangeCallback) {

	for wresp := range rch {
		for _, ev := range wresp.Events {
			onChangeCallback(convertUIntArrToString(ev.Kv.Key), convertUIntArrToString(ev.Kv.Value))
		}
	}

}

// NewEtcdCli constructs a new EtcdClient and checks etcd port's availability
func NewEtcdCli(conf Config) (etcdCli *EtcdCli, err error) {

	var cfg clientv3.Config
	endpoint := strings.Split(conf.Endpoint[0], ":")

	portUp := util.CheckPortAvailability(endpoint[0], endpoint[1])
	if !portUp {
		glog.Error("etcd service port not up")
		return nil, errors.New("etcd service port not up")
	}

	devmode := false
	if conf.CertFile == "" && conf.KeyFile == "" && conf.TrustFile == "" {
		devmode = true
	}

	if devmode {
		cfg = clientv3.Config{
			Endpoints: conf.Endpoint,
		}
	} else {
		tlsInfo := transport.TLSInfo{
			CertFile:      conf.CertFile,
			KeyFile:       conf.KeyFile,
			TrustedCAFile: conf.TrustFile,
		}

		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}

		cfg = clientv3.Config{
			Endpoints: conf.Endpoint,
			TLS:       tlsConfig,
		}
	}

	etcdcli, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &EtcdCli{etcd: etcdcli}, err
}

// GetConfig gets the value of a key from Etcd
func (etcdClient *EtcdCli) GetConfig(key string) (string, error) {

	response, err := etcdClient.etcd.Get(context.TODO(), key)
	if err != nil {
		return "", err
	}
	responseString := ""
	for _, ev := range response.Kvs {
		responseString = convertUIntArrToString(ev.Value)
	}
	return responseString, err
}

// RegisterDirWatch registers to a callback and keeps a watch on the prefix of a specified key
func (etcdClient *EtcdCli) RegisterDirWatch(key string, onChangeCallback OnChangeCallback) {

	glog.Infoln("watching on Dir:", key)
	rch := etcdClient.etcd.Watch(context.TODO(), key, clientv3.WithPrefix())
	go RegisterWatchOnKey(rch, onChangeCallback)

}

// RegisterKeyWatch registers to a callback and keeps a watch on a specified key
func (etcdClient *EtcdCli) RegisterKeyWatch(key string, onChangeCallback OnChangeCallback) {

	glog.Infoln("Watching on key:", key)
	rch := etcdClient.etcd.Watch(context.TODO(), key)
	go RegisterWatchOnKey(rch, onChangeCallback)

}

//convertUIntArrToString function to convert unit8[] to string
func convertUIntArrToString(uintarray []uint8) string {
	b := make([]byte, len(uintarray))
	for i, v := range uintarray {
		b[i] = byte(v)
	}
	return string(b)
}
