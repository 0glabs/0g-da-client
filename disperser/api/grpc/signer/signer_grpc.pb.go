// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.5
// source: signer/signer.proto

package signer

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SignerClient is the client API for Signer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SignerClient interface {
	BatchSign(ctx context.Context, in *BatchSignRequest, opts ...grpc.CallOption) (*BatchSignReply, error)
}

type signerClient struct {
	cc grpc.ClientConnInterface
}

func NewSignerClient(cc grpc.ClientConnInterface) SignerClient {
	return &signerClient{cc}
}

func (c *signerClient) BatchSign(ctx context.Context, in *BatchSignRequest, opts ...grpc.CallOption) (*BatchSignReply, error) {
	out := new(BatchSignReply)
	err := c.cc.Invoke(ctx, "/signer.Signer/BatchSign", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SignerServer is the server API for Signer service.
// All implementations must embed UnimplementedSignerServer
// for forward compatibility
type SignerServer interface {
	BatchSign(context.Context, *BatchSignRequest) (*BatchSignReply, error)
	mustEmbedUnimplementedSignerServer()
}

// UnimplementedSignerServer must be embedded to have forward compatible implementations.
type UnimplementedSignerServer struct {
}

func (UnimplementedSignerServer) BatchSign(context.Context, *BatchSignRequest) (*BatchSignReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchSign not implemented")
}
func (UnimplementedSignerServer) mustEmbedUnimplementedSignerServer() {}

// UnsafeSignerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SignerServer will
// result in compilation errors.
type UnsafeSignerServer interface {
	mustEmbedUnimplementedSignerServer()
}

func RegisterSignerServer(s grpc.ServiceRegistrar, srv SignerServer) {
	s.RegisterService(&Signer_ServiceDesc, srv)
}

func _Signer_BatchSign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchSignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SignerServer).BatchSign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/signer.Signer/BatchSign",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SignerServer).BatchSign(ctx, req.(*BatchSignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Signer_ServiceDesc is the grpc.ServiceDesc for Signer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Signer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "signer.Signer",
	HandlerType: (*SignerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BatchSign",
			Handler:    _Signer_BatchSign_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "signer/signer.proto",
}
