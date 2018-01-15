// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/antha-lang/antha/driver/antha_platereader_v1/platereader.proto

/*
Package antha_platereader_v1 is a generated protocol buffer package.

It is generated from these files:
	github.com/antha-lang/antha/driver/antha_platereader_v1/platereader.proto

It has these top-level messages:
	BoolReply
	ProtocolRunRequest
*/
package antha_platereader_v1

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

type BoolReply struct {
	Result bool `protobuf:"varint,1,opt,name=result" json:"result,omitempty"`
}

func (m *BoolReply) Reset()                    { *m = BoolReply{} }
func (m *BoolReply) String() string            { return proto.CompactTextString(m) }
func (*BoolReply) ProtoMessage()               {}
func (*BoolReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *BoolReply) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

type ProtocolRunRequest struct {
	ProtocolName string `protobuf:"bytes,1,opt,name=ProtocolName" json:"ProtocolName,omitempty"`
	PlateID      string `protobuf:"bytes,2,opt,name=PlateID" json:"PlateID,omitempty"`
}

func (m *ProtocolRunRequest) Reset()                    { *m = ProtocolRunRequest{} }
func (m *ProtocolRunRequest) String() string            { return proto.CompactTextString(m) }
func (*ProtocolRunRequest) ProtoMessage()               {}
func (*ProtocolRunRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ProtocolRunRequest) GetProtocolName() string {
	if m != nil {
		return m.ProtocolName
	}
	return ""
}

func (m *ProtocolRunRequest) GetPlateID() string {
	if m != nil {
		return m.PlateID
	}
	return ""
}

func init() {
	proto.RegisterType((*BoolReply)(nil), "antha.platereader.v1.BoolReply")
	proto.RegisterType((*ProtocolRunRequest)(nil), "antha.platereader.v1.ProtocolRunRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for PlateReader service

type PlateReaderClient interface {
	PRRunProtocolByName(ctx context.Context, in *ProtocolRunRequest, opts ...grpc.CallOption) (*BoolReply, error)
}

type plateReaderClient struct {
	cc *grpc.ClientConn
}

func NewPlateReaderClient(cc *grpc.ClientConn) PlateReaderClient {
	return &plateReaderClient{cc}
}

func (c *plateReaderClient) PRRunProtocolByName(ctx context.Context, in *ProtocolRunRequest, opts ...grpc.CallOption) (*BoolReply, error) {
	out := new(BoolReply)
	err := grpc.Invoke(ctx, "/antha.platereader.v1.PlateReader/PRRunProtocolByName", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for PlateReader service

type PlateReaderServer interface {
	PRRunProtocolByName(context.Context, *ProtocolRunRequest) (*BoolReply, error)
}

func RegisterPlateReaderServer(s *grpc.Server, srv PlateReaderServer) {
	s.RegisterService(&_PlateReader_serviceDesc, srv)
}

func _PlateReader_PRRunProtocolByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProtocolRunRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlateReaderServer).PRRunProtocolByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/antha.platereader.v1.PlateReader/PRRunProtocolByName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlateReaderServer).PRRunProtocolByName(ctx, req.(*ProtocolRunRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PlateReader_serviceDesc = grpc.ServiceDesc{
	ServiceName: "antha.platereader.v1.PlateReader",
	HandlerType: (*PlateReaderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PRRunProtocolByName",
			Handler:    _PlateReader_PRRunProtocolByName_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/antha-lang/antha/driver/antha_platereader_v1/platereader.proto",
}

func init() {
	proto.RegisterFile("github.com/antha-lang/antha/driver/antha_platereader_v1/platereader.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 219 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xf2, 0x4c, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xcc, 0x2b, 0xc9, 0x48, 0xd4, 0xcd, 0x49, 0xcc,
	0x4b, 0x87, 0x30, 0xf5, 0x53, 0x8a, 0x32, 0xcb, 0x52, 0x8b, 0x20, 0x9c, 0xf8, 0x82, 0x9c, 0xc4,
	0x92, 0xd4, 0xa2, 0xd4, 0xc4, 0x94, 0xd4, 0xa2, 0xf8, 0x32, 0x43, 0x7d, 0x24, 0xae, 0x5e, 0x41,
	0x51, 0x7e, 0x49, 0xbe, 0x90, 0x08, 0x58, 0x9d, 0x1e, 0xb2, 0x44, 0x99, 0xa1, 0x92, 0x32, 0x17,
	0xa7, 0x53, 0x7e, 0x7e, 0x4e, 0x50, 0x6a, 0x41, 0x4e, 0xa5, 0x90, 0x18, 0x17, 0x5b, 0x51, 0x6a,
	0x71, 0x69, 0x4e, 0x89, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x47, 0x10, 0x94, 0xa7, 0x14, 0xc4, 0x25,
	0x14, 0x00, 0x32, 0x23, 0x39, 0x3f, 0x27, 0xa8, 0x34, 0x2f, 0x28, 0xb5, 0xb0, 0x34, 0xb5, 0xb8,
	0x44, 0x48, 0x89, 0x8b, 0x07, 0x26, 0xea, 0x97, 0x98, 0x9b, 0x0a, 0xd6, 0xc3, 0x19, 0x84, 0x22,
	0x26, 0x24, 0xc1, 0xc5, 0x1e, 0x00, 0xb2, 0xd0, 0xd3, 0x45, 0x82, 0x09, 0x2c, 0x0d, 0xe3, 0x1a,
	0x15, 0x72, 0x71, 0x83, 0x99, 0x41, 0x60, 0xa7, 0x08, 0x25, 0x71, 0x09, 0x07, 0x04, 0x05, 0x95,
	0xe6, 0xc1, 0x74, 0x3b, 0x55, 0x82, 0xf5, 0x6b, 0xe8, 0x61, 0x73, 0xb5, 0x1e, 0xa6, 0x6b, 0xa4,
	0xe4, 0xb1, 0xab, 0x84, 0x7b, 0x4e, 0x89, 0x21, 0x89, 0x0d, 0x1c, 0x10, 0xc6, 0x80, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x69, 0x06, 0x4c, 0x6e, 0x55, 0x01, 0x00, 0x00,
}