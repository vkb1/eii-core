/*
Copyright (c) 2019 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package msgbus

import (
	types "EISMessageBus/pkg/types"
	"reflect"
	"testing"
	"time"
)

func initMsgBus() (*MsgbusContext, error) {
	m := map[string]interface{}{"type": "zmq_ipc", "socket_dir": "/opt/socks"}
	ctx, err := NewMsgbusContext(m)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func TestMsgContextInit(t *testing.T) {
	ctx, err := initMsgBus()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer ctx.Destroy()
}

func TestPubSub(t *testing.T) {
	ctx, err := initMsgBus()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer ctx.Destroy()

	pubCtx, err := ctx.NewPublisher("test")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	defer ctx.DestroyPublisher(pubCtx)

	subCtx, err := ctx.NewSubscriber("test")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	defer ctx.DestroyRecvCtx(subCtx)

	time.Sleep(1000 * time.Millisecond)

	m := map[string]interface{}{"str": "hello", "int": 2.0, "float": 55.5, "bool": true}
	env := types.NewMsgEnvelope(m, nil)

	err = ctx.Publish(pubCtx, env)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	received, err := ctx.ReceiveWait(subCtx)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	eq := reflect.DeepEqual(received, env)
	if !eq {
		t.Errorf("Messages not equal:\n\tSENT: %v\n\tRECEIVED: %v\n", env, received)
		return
	}
}

func TestReceiveNoWait(t *testing.T) {
	ctx, err := initMsgBus()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer ctx.Destroy()

	subCtx, err := ctx.NewSubscriber("asdf")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer ctx.DestroyRecvCtx(subCtx)

	msg, err := ctx.ReceiveNoWait(subCtx)
	if msg != nil || err != nil {
		t.Errorf("Both should be nil:\n\tmsg: %v\n\terr: %v\n", msg, err)
		return
	}
}

func TestReceiveTimedWait(t *testing.T) {
	ctx, err := initMsgBus()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer ctx.Destroy()

	subCtx, err := ctx.NewSubscriber("asdf")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer ctx.DestroyRecvCtx(subCtx)

	msg, err := ctx.ReceiveTimedWait(subCtx, 100)
	if msg != nil || err != nil {
		t.Errorf("Both should be nil:\n\tmsg: %v\n\terr: %v\n", msg, err)
		return
	}
}

func TestRequestResponse(t *testing.T) {
	ctx, err := initMsgBus()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer ctx.Destroy()

	serviceCtx, err := ctx.NewService("unittest-service")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer ctx.DestroyRecvCtx(serviceCtx)

	reqCtx, err := ctx.GetService("unittest-service")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer ctx.DestroyRecvCtx(reqCtx)

	// Send request
	req := types.NewMsgEnvelope(map[string]interface{}{"hello": "world"}, nil)

	err = ctx.Request(reqCtx, req)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	// Receive request
	receivedReq, err := ctx.ReceiveTimedWait(serviceCtx, 1000)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	eq := reflect.DeepEqual(req, receivedReq)
	if !eq {
		t.Errorf("Received request does not match request:\n\tSENT: %v\n\tRECEIVED: %v", req, receivedReq)
		return
	}

	// Send response
	err = ctx.Response(serviceCtx, receivedReq)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	// Receive response
	resp, err := ctx.ReceiveTimedWait(reqCtx, 1000)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	eq = reflect.DeepEqual(req, resp)
	if !eq {
		t.Errorf("Received response does not match original request:\n\tSENT: %v\n\tRECEIVED: %v",
			req, receivedReq)
		return
	}
}
