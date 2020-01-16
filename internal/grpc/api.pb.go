// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package grpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Response struct {
	Ok                   bool     `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

type AuthRequest struct {
	Login                string   `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Ip                   string   `protobuf:"bytes,3,opt,name=ip,proto3" json:"ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthRequest) Reset()         { *m = AuthRequest{} }
func (m *AuthRequest) String() string { return proto.CompactTextString(m) }
func (*AuthRequest) ProtoMessage()    {}
func (*AuthRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *AuthRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthRequest.Unmarshal(m, b)
}
func (m *AuthRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthRequest.Marshal(b, m, deterministic)
}
func (m *AuthRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthRequest.Merge(m, src)
}
func (m *AuthRequest) XXX_Size() int {
	return xxx_messageInfo_AuthRequest.Size(m)
}
func (m *AuthRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AuthRequest proto.InternalMessageInfo

func (m *AuthRequest) GetLogin() string {
	if m != nil {
		return m.Login
	}
	return ""
}

func (m *AuthRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *AuthRequest) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

type ClearBucketRequest struct {
	Login                string   `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Ip                   string   `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClearBucketRequest) Reset()         { *m = ClearBucketRequest{} }
func (m *ClearBucketRequest) String() string { return proto.CompactTextString(m) }
func (*ClearBucketRequest) ProtoMessage()    {}
func (*ClearBucketRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

func (m *ClearBucketRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClearBucketRequest.Unmarshal(m, b)
}
func (m *ClearBucketRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClearBucketRequest.Marshal(b, m, deterministic)
}
func (m *ClearBucketRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClearBucketRequest.Merge(m, src)
}
func (m *ClearBucketRequest) XXX_Size() int {
	return xxx_messageInfo_ClearBucketRequest.Size(m)
}
func (m *ClearBucketRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ClearBucketRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ClearBucketRequest proto.InternalMessageInfo

func (m *ClearBucketRequest) GetLogin() string {
	if m != nil {
		return m.Login
	}
	return ""
}

func (m *ClearBucketRequest) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

type IpListRequest struct {
	Ip                   string   `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IpListRequest) Reset()         { *m = IpListRequest{} }
func (m *IpListRequest) String() string { return proto.CompactTextString(m) }
func (*IpListRequest) ProtoMessage()    {}
func (*IpListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}

func (m *IpListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IpListRequest.Unmarshal(m, b)
}
func (m *IpListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IpListRequest.Marshal(b, m, deterministic)
}
func (m *IpListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IpListRequest.Merge(m, src)
}
func (m *IpListRequest) XXX_Size() int {
	return xxx_messageInfo_IpListRequest.Size(m)
}
func (m *IpListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IpListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IpListRequest proto.InternalMessageInfo

func (m *IpListRequest) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func init() {
	proto.RegisterType((*Response)(nil), "grpc.Response")
	proto.RegisterType((*AuthRequest)(nil), "grpc.AuthRequest")
	proto.RegisterType((*ClearBucketRequest)(nil), "grpc.ClearBucketRequest")
	proto.RegisterType((*IpListRequest)(nil), "grpc.IpListRequest")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0xcd, 0xb6, 0x4a, 0x3a, 0xc5, 0xa2, 0xa3, 0x42, 0xc8, 0x45, 0xc9, 0x49, 0x10, 0x72,
	0xd0, 0x53, 0x3d, 0x08, 0xa9, 0x22, 0x04, 0x04, 0x21, 0x17, 0xcf, 0x31, 0x19, 0xda, 0x25, 0x31,
	0x3b, 0xee, 0x6e, 0xf0, 0x07, 0xfb, 0x47, 0x24, 0x59, 0x1a, 0x2b, 0x16, 0xa1, 0x78, 0x9c, 0x7d,
	0xef, 0x9b, 0x37, 0x0f, 0x16, 0x26, 0x39, 0xcb, 0x98, 0xb5, 0xb2, 0x0a, 0xc7, 0x4b, 0xcd, 0x45,
	0x14, 0x82, 0x9f, 0x91, 0x61, 0xd5, 0x18, 0xc2, 0x19, 0x08, 0x55, 0x05, 0xde, 0x85, 0x77, 0xe9,
	0x67, 0x42, 0x55, 0xd1, 0x33, 0x4c, 0x93, 0xd6, 0xae, 0x32, 0x7a, 0x6f, 0xc9, 0x58, 0x3c, 0x85,
	0xfd, 0x5a, 0x2d, 0x65, 0xd3, 0x3b, 0x26, 0x99, 0x1b, 0x30, 0x04, 0x9f, 0x73, 0x63, 0x3e, 0x94,
	0x2e, 0x03, 0xd1, 0x0b, 0xc3, 0xdc, 0x2d, 0x94, 0x1c, 0x8c, 0xfa, 0x57, 0x21, 0x39, 0xba, 0x05,
	0xbc, 0xaf, 0x29, 0xd7, 0x8b, 0xb6, 0xa8, 0xc8, 0xfe, 0xbd, 0xd7, 0xb1, 0x62, 0x60, 0xcf, 0xe1,
	0x30, 0xe5, 0x27, 0x69, 0x06, 0xcc, 0x19, 0xbc, 0xb5, 0xe1, 0xfa, 0x53, 0xc0, 0x28, 0x61, 0x89,
	0x73, 0x98, 0x6e, 0x84, 0x60, 0x10, 0x77, 0x3d, 0xe3, 0xdf, 0xb9, 0xe1, 0xcc, 0x29, 0xeb, 0xfa,
	0xd1, 0x1e, 0x5e, 0xc1, 0xb8, 0x2b, 0x8c, 0xc7, 0x4e, 0xd9, 0x28, 0xbf, 0xc5, 0x3c, 0x87, 0xa3,
	0xa4, 0x2c, 0x53, 0x4e, 0x9b, 0x45, 0x9d, 0x17, 0x55, 0x77, 0x1a, 0x9e, 0x38, 0xd7, 0x8f, 0x43,
	0xb7, 0xa0, 0x77, 0x70, 0xf6, 0x40, 0x35, 0x59, 0x4a, 0xf9, 0x51, 0xab, 0xb7, 0x9d, 0xf9, 0xef,
	0xe8, 0x97, 0x95, 0xb4, 0xf4, 0x8f, 0xe8, 0x5d, 0xf9, 0xd7, 0x83, 0xfe, 0xf3, 0xdc, 0x7c, 0x05,
	0x00, 0x00, 0xff, 0xff, 0x40, 0x94, 0xc2, 0xc0, 0x49, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ApiClient is the client API for Api service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ApiClient interface {
	ClearBucket(ctx context.Context, in *ClearBucketRequest, opts ...grpc.CallOption) (*Response, error)
	Auth(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*Response, error)
	AddIpInBlackList(ctx context.Context, in *IpListRequest, opts ...grpc.CallOption) (*Response, error)
	DeleteIpFromBlackList(ctx context.Context, in *IpListRequest, opts ...grpc.CallOption) (*Response, error)
	AddIpInWhiteList(ctx context.Context, in *IpListRequest, opts ...grpc.CallOption) (*Response, error)
	DeleteIpFromWhiteList(ctx context.Context, in *IpListRequest, opts ...grpc.CallOption) (*Response, error)
}

type apiClient struct {
	cc *grpc.ClientConn
}

func NewApiClient(cc *grpc.ClientConn) ApiClient {
	return &apiClient{cc}
}

func (c *apiClient) ClearBucket(ctx context.Context, in *ClearBucketRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.Api/ClearBucket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Auth(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.Api/Auth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) AddIpInBlackList(ctx context.Context, in *IpListRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.Api/AddIpInBlackList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) DeleteIpFromBlackList(ctx context.Context, in *IpListRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.Api/DeleteIpFromBlackList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) AddIpInWhiteList(ctx context.Context, in *IpListRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.Api/AddIpInWhiteList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) DeleteIpFromWhiteList(ctx context.Context, in *IpListRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/grpc.Api/DeleteIpFromWhiteList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApiServer is the server API for Api service.
type ApiServer interface {
	ClearBucket(context.Context, *ClearBucketRequest) (*Response, error)
	Auth(context.Context, *AuthRequest) (*Response, error)
	AddIpInBlackList(context.Context, *IpListRequest) (*Response, error)
	DeleteIpFromBlackList(context.Context, *IpListRequest) (*Response, error)
	AddIpInWhiteList(context.Context, *IpListRequest) (*Response, error)
	DeleteIpFromWhiteList(context.Context, *IpListRequest) (*Response, error)
}

// UnimplementedApiServer can be embedded to have forward compatible implementations.
type UnimplementedApiServer struct {
}

func (*UnimplementedApiServer) ClearBucket(ctx context.Context, req *ClearBucketRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearBucket not implemented")
}
func (*UnimplementedApiServer) Auth(ctx context.Context, req *AuthRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Auth not implemented")
}
func (*UnimplementedApiServer) AddIpInBlackList(ctx context.Context, req *IpListRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddIpInBlackList not implemented")
}
func (*UnimplementedApiServer) DeleteIpFromBlackList(ctx context.Context, req *IpListRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteIpFromBlackList not implemented")
}
func (*UnimplementedApiServer) AddIpInWhiteList(ctx context.Context, req *IpListRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddIpInWhiteList not implemented")
}
func (*UnimplementedApiServer) DeleteIpFromWhiteList(ctx context.Context, req *IpListRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteIpFromWhiteList not implemented")
}

func RegisterApiServer(s *grpc.Server, srv ApiServer) {
	s.RegisterService(&_Api_serviceDesc, srv)
}

func _Api_ClearBucket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearBucketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).ClearBucket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Api/ClearBucket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).ClearBucket(ctx, req.(*ClearBucketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Auth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Auth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Api/Auth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Auth(ctx, req.(*AuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_AddIpInBlackList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IpListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).AddIpInBlackList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Api/AddIpInBlackList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).AddIpInBlackList(ctx, req.(*IpListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_DeleteIpFromBlackList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IpListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).DeleteIpFromBlackList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Api/DeleteIpFromBlackList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).DeleteIpFromBlackList(ctx, req.(*IpListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_AddIpInWhiteList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IpListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).AddIpInWhiteList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Api/AddIpInWhiteList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).AddIpInWhiteList(ctx, req.(*IpListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_DeleteIpFromWhiteList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IpListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).DeleteIpFromWhiteList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Api/DeleteIpFromWhiteList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).DeleteIpFromWhiteList(ctx, req.(*IpListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Api_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Api",
	HandlerType: (*ApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClearBucket",
			Handler:    _Api_ClearBucket_Handler,
		},
		{
			MethodName: "Auth",
			Handler:    _Api_Auth_Handler,
		},
		{
			MethodName: "AddIpInBlackList",
			Handler:    _Api_AddIpInBlackList_Handler,
		},
		{
			MethodName: "DeleteIpFromBlackList",
			Handler:    _Api_DeleteIpFromBlackList_Handler,
		},
		{
			MethodName: "AddIpInWhiteList",
			Handler:    _Api_AddIpInWhiteList_Handler,
		},
		{
			MethodName: "DeleteIpFromWhiteList",
			Handler:    _Api_DeleteIpFromWhiteList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
