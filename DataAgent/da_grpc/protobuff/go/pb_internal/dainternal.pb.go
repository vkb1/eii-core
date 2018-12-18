// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dainternal.proto

package DataAgentInternal

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message for GetConfigInt
type ConfigIntReq struct {
	CfgType              string   `protobuf:"bytes,1,opt,name=cfgType" json:"cfgType,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigIntReq) Reset()         { *m = ConfigIntReq{} }
func (m *ConfigIntReq) String() string { return proto.CompactTextString(m) }
func (*ConfigIntReq) ProtoMessage()    {}
func (*ConfigIntReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_dainternal_1e8462ed3c216bbd, []int{0}
}
func (m *ConfigIntReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigIntReq.Unmarshal(m, b)
}
func (m *ConfigIntReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigIntReq.Marshal(b, m, deterministic)
}
func (dst *ConfigIntReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigIntReq.Merge(dst, src)
}
func (m *ConfigIntReq) XXX_Size() int {
	return xxx_messageInfo_ConfigIntReq.Size(m)
}
func (m *ConfigIntReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigIntReq.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigIntReq proto.InternalMessageInfo

func (m *ConfigIntReq) GetCfgType() string {
	if m != nil {
		return m.CfgType
	}
	return ""
}

// The response message for GetConfigInt
type ConfigIntResp struct {
	JsonMsg              string   `protobuf:"bytes,1,opt,name=jsonMsg" json:"jsonMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigIntResp) Reset()         { *m = ConfigIntResp{} }
func (m *ConfigIntResp) String() string { return proto.CompactTextString(m) }
func (*ConfigIntResp) ProtoMessage()    {}
func (*ConfigIntResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_dainternal_1e8462ed3c216bbd, []int{1}
}
func (m *ConfigIntResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigIntResp.Unmarshal(m, b)
}
func (m *ConfigIntResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigIntResp.Marshal(b, m, deterministic)
}
func (dst *ConfigIntResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigIntResp.Merge(dst, src)
}
func (m *ConfigIntResp) XXX_Size() int {
	return xxx_messageInfo_ConfigIntResp.Size(m)
}
func (m *ConfigIntResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigIntResp.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigIntResp proto.InternalMessageInfo

func (m *ConfigIntResp) GetJsonMsg() string {
	if m != nil {
		return m.JsonMsg
	}
	return ""
}

func init() {
	proto.RegisterType((*ConfigIntReq)(nil), "DataAgentInternal.ConfigIntReq")
	proto.RegisterType((*ConfigIntResp)(nil), "DataAgentInternal.ConfigIntResp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DainternalClient is the client API for Dainternal service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DainternalClient interface {
	// **********Internal Interfaces***************
	// GetConfigInt internal interface
	GetConfigInt(ctx context.Context, in *ConfigIntReq, opts ...grpc.CallOption) (*ConfigIntResp, error)
}

type dainternalClient struct {
	cc *grpc.ClientConn
}

func NewDainternalClient(cc *grpc.ClientConn) DainternalClient {
	return &dainternalClient{cc}
}

func (c *dainternalClient) GetConfigInt(ctx context.Context, in *ConfigIntReq, opts ...grpc.CallOption) (*ConfigIntResp, error) {
	out := new(ConfigIntResp)
	err := c.cc.Invoke(ctx, "/DataAgentInternal.dainternal/GetConfigInt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DainternalServer is the server API for Dainternal service.
type DainternalServer interface {
	// **********Internal Interfaces***************
	// GetConfigInt internal interface
	GetConfigInt(context.Context, *ConfigIntReq) (*ConfigIntResp, error)
}

func RegisterDainternalServer(s *grpc.Server, srv DainternalServer) {
	s.RegisterService(&_Dainternal_serviceDesc, srv)
}

func _Dainternal_GetConfigInt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigIntReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DainternalServer).GetConfigInt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DataAgentInternal.dainternal/GetConfigInt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DainternalServer).GetConfigInt(ctx, req.(*ConfigIntReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Dainternal_serviceDesc = grpc.ServiceDesc{
	ServiceName: "DataAgentInternal.dainternal",
	HandlerType: (*DainternalServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetConfigInt",
			Handler:    _Dainternal_GetConfigInt_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dainternal.proto",
}

func init() { proto.RegisterFile("dainternal.proto", fileDescriptor_dainternal_1e8462ed3c216bbd) }

var fileDescriptor_dainternal_1e8462ed3c216bbd = []byte{
	// 150 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x48, 0x49, 0xcc, 0xcc,
	0x2b, 0x49, 0x2d, 0xca, 0x4b, 0xcc, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x74, 0x49,
	0x2c, 0x49, 0x74, 0x4c, 0x4f, 0xcd, 0x2b, 0xf1, 0x84, 0x4a, 0x28, 0x69, 0x70, 0xf1, 0x38, 0xe7,
	0xe7, 0xa5, 0x65, 0xa6, 0x7b, 0xe6, 0x95, 0x04, 0xa5, 0x16, 0x0a, 0x49, 0x70, 0xb1, 0x27, 0xa7,
	0xa5, 0x87, 0x54, 0x16, 0xa4, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xb8, 0x4a, 0x9a,
	0x5c, 0xbc, 0x48, 0x2a, 0x8b, 0x0b, 0x40, 0x4a, 0xb3, 0x8a, 0xf3, 0xf3, 0x7c, 0x8b, 0xd3, 0x61,
	0x4a, 0xa1, 0x5c, 0xa3, 0x44, 0x2e, 0x2e, 0x84, 0xdd, 0x42, 0xc1, 0x5c, 0x3c, 0xee, 0xa9, 0x25,
	0x70, 0xbd, 0x42, 0xf2, 0x7a, 0x18, 0xce, 0xd0, 0x43, 0x76, 0x83, 0x94, 0x02, 0x7e, 0x05, 0xc5,
	0x05, 0x4a, 0x0c, 0x49, 0x6c, 0x60, 0x1f, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x39, 0x0a,
	0xa0, 0x3b, 0xe5, 0x00, 0x00, 0x00,
}