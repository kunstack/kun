// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: peer.proto

package v1

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

// PeerServiceClient is the client API for PeerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PeerServiceClient interface {
	// SignUpUser 注册账号
	WatchUpstream(ctx context.Context, in *WatchUpstreamRequest, opts ...grpc.CallOption) (PeerService_WatchUpstreamClient, error)
}

type peerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPeerServiceClient(cc grpc.ClientConnInterface) PeerServiceClient {
	return &peerServiceClient{cc}
}

func (c *peerServiceClient) WatchUpstream(ctx context.Context, in *WatchUpstreamRequest, opts ...grpc.CallOption) (PeerService_WatchUpstreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &PeerService_ServiceDesc.Streams[0], "/compass.api.v1.PeerService/WatchUpstream", opts...)
	if err != nil {
		return nil, err
	}
	x := &peerServiceWatchUpstreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PeerService_WatchUpstreamClient interface {
	Recv() (*WatchUpstreamResponse, error)
	grpc.ClientStream
}

type peerServiceWatchUpstreamClient struct {
	grpc.ClientStream
}

func (x *peerServiceWatchUpstreamClient) Recv() (*WatchUpstreamResponse, error) {
	m := new(WatchUpstreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PeerServiceServer is the server API for PeerService service.
// All implementations must embed UnimplementedPeerServiceServer
// for forward compatibility
type PeerServiceServer interface {
	// SignUpUser 注册账号
	WatchUpstream(*WatchUpstreamRequest, PeerService_WatchUpstreamServer) error
	mustEmbedUnimplementedPeerServiceServer()
}

// UnimplementedPeerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPeerServiceServer struct {
}

func (UnimplementedPeerServiceServer) WatchUpstream(*WatchUpstreamRequest, PeerService_WatchUpstreamServer) error {
	return status.Errorf(codes.Unimplemented, "method WatchUpstream not implemented")
}
func (UnimplementedPeerServiceServer) mustEmbedUnimplementedPeerServiceServer() {}

// UnsafePeerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PeerServiceServer will
// result in compilation errors.
type UnsafePeerServiceServer interface {
	mustEmbedUnimplementedPeerServiceServer()
}

func RegisterPeerServiceServer(s grpc.ServiceRegistrar, srv PeerServiceServer) {
	s.RegisterService(&PeerService_ServiceDesc, srv)
}

func _PeerService_WatchUpstream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchUpstreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PeerServiceServer).WatchUpstream(m, &peerServiceWatchUpstreamServer{stream})
}

type PeerService_WatchUpstreamServer interface {
	Send(*WatchUpstreamResponse) error
	grpc.ServerStream
}

type peerServiceWatchUpstreamServer struct {
	grpc.ServerStream
}

func (x *peerServiceWatchUpstreamServer) Send(m *WatchUpstreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// PeerService_ServiceDesc is the grpc.ServiceDesc for PeerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PeerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "compass.api.v1.PeerService",
	HandlerType: (*PeerServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchUpstream",
			Handler:       _PeerService_WatchUpstream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "peer.proto",
}
