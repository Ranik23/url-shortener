// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: service.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	URLShortener_ShortenURL_FullMethodName     = "/urlshortener.URLShortener/ShortenURL"
	URLShortener_GetOriginalURL_FullMethodName = "/urlshortener.URLShortener/GetOriginalURL"
	URLShortener_GetStats_FullMethodName       = "/urlshortener.URLShortener/GetStats"
	URLShortener_DeleteURL_FullMethodName      = "/urlshortener.URLShortener/DeleteURL"
)

// URLShortenerClient is the client API for URLShortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLShortenerClient interface {
	ShortenURL(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error)
	GetOriginalURL(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	GetStats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsResponse, error)
	DeleteURL(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
}

type uRLShortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewURLShortenerClient(cc grpc.ClientConnInterface) URLShortenerClient {
	return &uRLShortenerClient{cc}
}

func (c *uRLShortenerClient) ShortenURL(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShortenResponse)
	err := c.cc.Invoke(ctx, URLShortener_ShortenURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) GetOriginalURL(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, URLShortener_GetOriginalURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) GetStats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, URLShortener_GetStats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) DeleteURL(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, URLShortener_DeleteURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLShortenerServer is the server API for URLShortener service.
// All implementations must embed UnimplementedURLShortenerServer
// for forward compatibility.
type URLShortenerServer interface {
	ShortenURL(context.Context, *ShortenRequest) (*ShortenResponse, error)
	GetOriginalURL(context.Context, *GetRequest) (*GetResponse, error)
	GetStats(context.Context, *StatsRequest) (*StatsResponse, error)
	DeleteURL(context.Context, *DeleteRequest) (*DeleteResponse, error)
	mustEmbedUnimplementedURLShortenerServer()
}

// UnimplementedURLShortenerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedURLShortenerServer struct{}

func (UnimplementedURLShortenerServer) ShortenURL(context.Context, *ShortenRequest) (*ShortenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenURL not implemented")
}
func (UnimplementedURLShortenerServer) GetOriginalURL(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOriginalURL not implemented")
}
func (UnimplementedURLShortenerServer) GetStats(context.Context, *StatsRequest) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStats not implemented")
}
func (UnimplementedURLShortenerServer) DeleteURL(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteURL not implemented")
}
func (UnimplementedURLShortenerServer) mustEmbedUnimplementedURLShortenerServer() {}
func (UnimplementedURLShortenerServer) testEmbeddedByValue()                      {}

// UnsafeURLShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLShortenerServer will
// result in compilation errors.
type UnsafeURLShortenerServer interface {
	mustEmbedUnimplementedURLShortenerServer()
}

func RegisterURLShortenerServer(s grpc.ServiceRegistrar, srv URLShortenerServer) {
	// If the following call pancis, it indicates UnimplementedURLShortenerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&URLShortener_ServiceDesc, srv)
}

func _URLShortener_ShortenURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).ShortenURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_ShortenURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).ShortenURL(ctx, req.(*ShortenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetOriginalURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetOriginalURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_GetOriginalURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetOriginalURL(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_GetStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetStats(ctx, req.(*StatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_DeleteURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).DeleteURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_DeleteURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).DeleteURL(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// URLShortener_ServiceDesc is the grpc.ServiceDesc for URLShortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLShortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "urlshortener.URLShortener",
	HandlerType: (*URLShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShortenURL",
			Handler:    _URLShortener_ShortenURL_Handler,
		},
		{
			MethodName: "GetOriginalURL",
			Handler:    _URLShortener_GetOriginalURL_Handler,
		},
		{
			MethodName: "GetStats",
			Handler:    _URLShortener_GetStats_Handler,
		},
		{
			MethodName: "DeleteURL",
			Handler:    _URLShortener_DeleteURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
