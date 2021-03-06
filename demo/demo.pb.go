// Code generated by protoc-gen-go. DO NOT EDIT.
// source: demo.proto

/*
Package demo is a generated protocol buffer package.

It is generated from these files:
	demo.proto

It has these top-level messages:
	EchoReq
	EchoResp
*/
package demo

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

type EchoReq struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *EchoReq) Reset()                    { *m = EchoReq{} }
func (m *EchoReq) String() string            { return proto.CompactTextString(m) }
func (*EchoReq) ProtoMessage()               {}
func (*EchoReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *EchoReq) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type EchoResp struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}

func (m *EchoResp) Reset()                    { *m = EchoResp{} }
func (m *EchoResp) String() string            { return proto.CompactTextString(m) }
func (*EchoResp) ProtoMessage()               {}
func (*EchoResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *EchoResp) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func init() {
	proto.RegisterType((*EchoReq)(nil), "demo.EchoReq")
	proto.RegisterType((*EchoResp)(nil), "demo.EchoResp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Demo service

type DemoClient interface {
	Echo(ctx context.Context, in *EchoReq, opts ...grpc.CallOption) (*EchoResp, error)
}

type demoClient struct {
	cc *grpc.ClientConn
}

func NewDemoClient(cc *grpc.ClientConn) DemoClient {
	return &demoClient{cc}
}

func (c *demoClient) Echo(ctx context.Context, in *EchoReq, opts ...grpc.CallOption) (*EchoResp, error) {
	out := new(EchoResp)
	err := grpc.Invoke(ctx, "/demo.Demo/Echo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Demo service

type DemoServer interface {
	Echo(context.Context, *EchoReq) (*EchoResp, error)
}

func RegisterDemoServer(s *grpc.Server, srv DemoServer) {
	s.RegisterService(&_Demo_serviceDesc, srv)
}

func _Demo_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DemoServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/demo.Demo/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DemoServer).Echo(ctx, req.(*EchoReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Demo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "demo.Demo",
	HandlerType: (*DemoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _Demo_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "demo.proto",
}

func init() { proto.RegisterFile("demo.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 127 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x49, 0xcd, 0xcd,
	0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0x95, 0x94, 0xb9, 0xd8, 0x5d, 0x93,
	0x33, 0xf2, 0x83, 0x52, 0x0b, 0x85, 0x24, 0xb8, 0xd8, 0x73, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53,
	0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x60, 0x5c, 0x25, 0x25, 0x2e, 0x0e, 0x88, 0xa2, 0xe2,
	0x02, 0x21, 0x31, 0x2e, 0xb6, 0xa2, 0xd4, 0xe2, 0xd2, 0x9c, 0x12, 0xa8, 0x22, 0x28, 0xcf, 0x48,
	0x9f, 0x8b, 0xc5, 0x25, 0x35, 0x37, 0x5f, 0x48, 0x9d, 0x8b, 0x05, 0xa4, 0x56, 0x88, 0x57, 0x0f,
	0x6c, 0x17, 0xd4, 0x70, 0x29, 0x3e, 0x64, 0x6e, 0x71, 0x81, 0x12, 0x43, 0x12, 0x1b, 0xd8, 0x19,
	0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd3, 0x4f, 0x8d, 0xcd, 0x94, 0x00, 0x00, 0x00,
}
