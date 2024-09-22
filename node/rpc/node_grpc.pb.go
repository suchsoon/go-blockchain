// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.1
// source: node.proto

package rpc

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
	Node_PeerDiscover_FullMethodName = "/Node/PeerDiscover"
	Node_EventStream_FullMethodName  = "/Node/EventStream"
)

// NodeClient is the client API for Node service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeClient interface {
	PeerDiscover(ctx context.Context, in *PeerDiscoverReq, opts ...grpc.CallOption) (*PeerDiscoverRes, error)
	EventStream(ctx context.Context, in *EventStreamReq, opts ...grpc.CallOption) (grpc.ServerStreamingClient[EventStreamRes], error)
}

type nodeClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeClient(cc grpc.ClientConnInterface) NodeClient {
	return &nodeClient{cc}
}

func (c *nodeClient) PeerDiscover(ctx context.Context, in *PeerDiscoverReq, opts ...grpc.CallOption) (*PeerDiscoverRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PeerDiscoverRes)
	err := c.cc.Invoke(ctx, Node_PeerDiscover_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) EventStream(ctx context.Context, in *EventStreamReq, opts ...grpc.CallOption) (grpc.ServerStreamingClient[EventStreamRes], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Node_ServiceDesc.Streams[0], Node_EventStream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[EventStreamReq, EventStreamRes]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Node_EventStreamClient = grpc.ServerStreamingClient[EventStreamRes]

// NodeServer is the server API for Node service.
// All implementations must embed UnimplementedNodeServer
// for forward compatibility.
type NodeServer interface {
	PeerDiscover(context.Context, *PeerDiscoverReq) (*PeerDiscoverRes, error)
	EventStream(*EventStreamReq, grpc.ServerStreamingServer[EventStreamRes]) error
	mustEmbedUnimplementedNodeServer()
}

// UnimplementedNodeServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNodeServer struct{}

func (UnimplementedNodeServer) PeerDiscover(context.Context, *PeerDiscoverReq) (*PeerDiscoverRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PeerDiscover not implemented")
}
func (UnimplementedNodeServer) EventStream(*EventStreamReq, grpc.ServerStreamingServer[EventStreamRes]) error {
	return status.Errorf(codes.Unimplemented, "method EventStream not implemented")
}
func (UnimplementedNodeServer) mustEmbedUnimplementedNodeServer() {}
func (UnimplementedNodeServer) testEmbeddedByValue()              {}

// UnsafeNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeServer will
// result in compilation errors.
type UnsafeNodeServer interface {
	mustEmbedUnimplementedNodeServer()
}

func RegisterNodeServer(s grpc.ServiceRegistrar, srv NodeServer) {
	// If the following call pancis, it indicates UnimplementedNodeServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Node_ServiceDesc, srv)
}

func _Node_PeerDiscover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PeerDiscoverReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).PeerDiscover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Node_PeerDiscover_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).PeerDiscover(ctx, req.(*PeerDiscoverReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_EventStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EventStreamReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NodeServer).EventStream(m, &grpc.GenericServerStream[EventStreamReq, EventStreamRes]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Node_EventStreamServer = grpc.ServerStreamingServer[EventStreamRes]

// Node_ServiceDesc is the grpc.ServiceDesc for Node service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Node_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Node",
	HandlerType: (*NodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PeerDiscover",
			Handler:    _Node_PeerDiscover_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "EventStream",
			Handler:       _Node_EventStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "node.proto",
}
