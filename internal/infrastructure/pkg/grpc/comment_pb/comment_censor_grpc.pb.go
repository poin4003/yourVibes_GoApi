// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: comment_censor.proto

package comment_pb

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
	CommentCensorService_CheckComment_FullMethodName = "/comment_censor.CommentCensorService/CheckComment"
)

// CommentCensorServiceClient is the client API for CommentCensorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommentCensorServiceClient interface {
	CheckComment(ctx context.Context, in *CommentCensorRequest, opts ...grpc.CallOption) (*CommentCensorResponse, error)
}

type commentCensorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCommentCensorServiceClient(cc grpc.ClientConnInterface) CommentCensorServiceClient {
	return &commentCensorServiceClient{cc}
}

func (c *commentCensorServiceClient) CheckComment(ctx context.Context, in *CommentCensorRequest, opts ...grpc.CallOption) (*CommentCensorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommentCensorResponse)
	err := c.cc.Invoke(ctx, CommentCensorService_CheckComment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommentCensorServiceServer is the server API for CommentCensorService service.
// All implementations must embed UnimplementedCommentCensorServiceServer
// for forward compatibility.
type CommentCensorServiceServer interface {
	CheckComment(context.Context, *CommentCensorRequest) (*CommentCensorResponse, error)
	mustEmbedUnimplementedCommentCensorServiceServer()
}

// UnimplementedCommentCensorServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCommentCensorServiceServer struct{}

func (UnimplementedCommentCensorServiceServer) CheckComment(context.Context, *CommentCensorRequest) (*CommentCensorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckComment not implemented")
}
func (UnimplementedCommentCensorServiceServer) mustEmbedUnimplementedCommentCensorServiceServer() {}
func (UnimplementedCommentCensorServiceServer) testEmbeddedByValue()                              {}

// UnsafeCommentCensorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommentCensorServiceServer will
// result in compilation errors.
type UnsafeCommentCensorServiceServer interface {
	mustEmbedUnimplementedCommentCensorServiceServer()
}

func RegisterCommentCensorServiceServer(s grpc.ServiceRegistrar, srv CommentCensorServiceServer) {
	// If the following call pancis, it indicates UnimplementedCommentCensorServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CommentCensorService_ServiceDesc, srv)
}

func _CommentCensorService_CheckComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommentCensorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentCensorServiceServer).CheckComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CommentCensorService_CheckComment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentCensorServiceServer).CheckComment(ctx, req.(*CommentCensorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CommentCensorService_ServiceDesc is the grpc.ServiceDesc for CommentCensorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommentCensorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "comment_censor.CommentCensorService",
	HandlerType: (*CommentCensorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckComment",
			Handler:    _CommentCensorService_CheckComment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "comment_censor.proto",
}
