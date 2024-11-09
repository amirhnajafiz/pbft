// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: pbft.proto

package pbft

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	PBFT_Request_FullMethodName     = "/pbft.PBFT/Request"
	PBFT_PrePrepare_FullMethodName  = "/pbft.PBFT/PrePrepare"
	PBFT_PrePrepared_FullMethodName = "/pbft.PBFT/PrePrepared"
	PBFT_Prepare_FullMethodName     = "/pbft.PBFT/Prepare"
	PBFT_Prepared_FullMethodName    = "/pbft.PBFT/Prepared"
	PBFT_Commit_FullMethodName      = "/pbft.PBFT/Commit"
	PBFT_ViewChange_FullMethodName  = "/pbft.PBFT/ViewChange"
	PBFT_NewView_FullMethodName     = "/pbft.PBFT/NewView"
	PBFT_Checkpoint_FullMethodName  = "/pbft.PBFT/Checkpoint"
	PBFT_PrintLog_FullMethodName    = "/pbft.PBFT/PrintLog"
	PBFT_PrintDB_FullMethodName     = "/pbft.PBFT/PrintDB"
	PBFT_PrintStatus_FullMethodName = "/pbft.PBFT/PrintStatus"
	PBFT_PrintView_FullMethodName   = "/pbft.PBFT/PrintView"
)

// PBFTClient is the client API for PBFT service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// creating rpc services for transactions and pbft.
// this service is for handling internal node calls for performing pbft.
type PBFTClient interface {
	Request(ctx context.Context, in *RequestMsg, opts ...grpc.CallOption) (*emptypb.Empty, error)
	PrePrepare(ctx context.Context, in *PrePrepareMsg, opts ...grpc.CallOption) (*emptypb.Empty, error)
	PrePrepared(ctx context.Context, in *AckMsg, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Prepare(ctx context.Context, in *AckMsg, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Prepared(ctx context.Context, in *AckMsg, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Commit(ctx context.Context, in *AckMsg, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ViewChange(ctx context.Context, in *ViewChangeMsg, opts ...grpc.CallOption) (*emptypb.Empty, error)
	NewView(ctx context.Context, in *NewViewMsg, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Checkpoint(ctx context.Context, in *CheckpointMsg, opts ...grpc.CallOption) (*emptypb.Empty, error)
	PrintLog(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[LogRsp], error)
	PrintDB(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[RequestMsg], error)
	PrintStatus(ctx context.Context, in *StatusMsg, opts ...grpc.CallOption) (*StatusRsp, error)
	PrintView(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ViewRsp], error)
}

type pBFTClient struct {
	cc grpc.ClientConnInterface
}

func NewPBFTClient(cc grpc.ClientConnInterface) PBFTClient {
	return &pBFTClient{cc}
}

func (c *pBFTClient) Request(ctx context.Context, in *RequestMsg, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PBFT_Request_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pBFTClient) PrePrepare(ctx context.Context, in *PrePrepareMsg, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PBFT_PrePrepare_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pBFTClient) PrePrepared(ctx context.Context, in *AckMsg, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PBFT_PrePrepared_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pBFTClient) Prepare(ctx context.Context, in *AckMsg, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PBFT_Prepare_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pBFTClient) Prepared(ctx context.Context, in *AckMsg, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PBFT_Prepared_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pBFTClient) Commit(ctx context.Context, in *AckMsg, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PBFT_Commit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pBFTClient) ViewChange(ctx context.Context, in *ViewChangeMsg, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PBFT_ViewChange_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pBFTClient) NewView(ctx context.Context, in *NewViewMsg, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PBFT_NewView_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pBFTClient) Checkpoint(ctx context.Context, in *CheckpointMsg, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PBFT_Checkpoint_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pBFTClient) PrintLog(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[LogRsp], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &PBFT_ServiceDesc.Streams[0], PBFT_PrintLog_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[emptypb.Empty, LogRsp]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type PBFT_PrintLogClient = grpc.ServerStreamingClient[LogRsp]

func (c *pBFTClient) PrintDB(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[RequestMsg], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &PBFT_ServiceDesc.Streams[1], PBFT_PrintDB_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[emptypb.Empty, RequestMsg]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type PBFT_PrintDBClient = grpc.ServerStreamingClient[RequestMsg]

func (c *pBFTClient) PrintStatus(ctx context.Context, in *StatusMsg, opts ...grpc.CallOption) (*StatusRsp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusRsp)
	err := c.cc.Invoke(ctx, PBFT_PrintStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pBFTClient) PrintView(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ViewRsp], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &PBFT_ServiceDesc.Streams[2], PBFT_PrintView_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[emptypb.Empty, ViewRsp]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type PBFT_PrintViewClient = grpc.ServerStreamingClient[ViewRsp]

// PBFTServer is the server API for PBFT service.
// All implementations must embed UnimplementedPBFTServer
// for forward compatibility.
//
// creating rpc services for transactions and pbft.
// this service is for handling internal node calls for performing pbft.
type PBFTServer interface {
	Request(context.Context, *RequestMsg) (*emptypb.Empty, error)
	PrePrepare(context.Context, *PrePrepareMsg) (*emptypb.Empty, error)
	PrePrepared(context.Context, *AckMsg) (*emptypb.Empty, error)
	Prepare(context.Context, *AckMsg) (*emptypb.Empty, error)
	Prepared(context.Context, *AckMsg) (*emptypb.Empty, error)
	Commit(context.Context, *AckMsg) (*emptypb.Empty, error)
	ViewChange(context.Context, *ViewChangeMsg) (*emptypb.Empty, error)
	NewView(context.Context, *NewViewMsg) (*emptypb.Empty, error)
	Checkpoint(context.Context, *CheckpointMsg) (*emptypb.Empty, error)
	PrintLog(*emptypb.Empty, grpc.ServerStreamingServer[LogRsp]) error
	PrintDB(*emptypb.Empty, grpc.ServerStreamingServer[RequestMsg]) error
	PrintStatus(context.Context, *StatusMsg) (*StatusRsp, error)
	PrintView(*emptypb.Empty, grpc.ServerStreamingServer[ViewRsp]) error
	mustEmbedUnimplementedPBFTServer()
}

// UnimplementedPBFTServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPBFTServer struct{}

func (UnimplementedPBFTServer) Request(context.Context, *RequestMsg) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Request not implemented")
}
func (UnimplementedPBFTServer) PrePrepare(context.Context, *PrePrepareMsg) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrePrepare not implemented")
}
func (UnimplementedPBFTServer) PrePrepared(context.Context, *AckMsg) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrePrepared not implemented")
}
func (UnimplementedPBFTServer) Prepare(context.Context, *AckMsg) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Prepare not implemented")
}
func (UnimplementedPBFTServer) Prepared(context.Context, *AckMsg) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Prepared not implemented")
}
func (UnimplementedPBFTServer) Commit(context.Context, *AckMsg) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Commit not implemented")
}
func (UnimplementedPBFTServer) ViewChange(context.Context, *ViewChangeMsg) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ViewChange not implemented")
}
func (UnimplementedPBFTServer) NewView(context.Context, *NewViewMsg) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewView not implemented")
}
func (UnimplementedPBFTServer) Checkpoint(context.Context, *CheckpointMsg) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Checkpoint not implemented")
}
func (UnimplementedPBFTServer) PrintLog(*emptypb.Empty, grpc.ServerStreamingServer[LogRsp]) error {
	return status.Errorf(codes.Unimplemented, "method PrintLog not implemented")
}
func (UnimplementedPBFTServer) PrintDB(*emptypb.Empty, grpc.ServerStreamingServer[RequestMsg]) error {
	return status.Errorf(codes.Unimplemented, "method PrintDB not implemented")
}
func (UnimplementedPBFTServer) PrintStatus(context.Context, *StatusMsg) (*StatusRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrintStatus not implemented")
}
func (UnimplementedPBFTServer) PrintView(*emptypb.Empty, grpc.ServerStreamingServer[ViewRsp]) error {
	return status.Errorf(codes.Unimplemented, "method PrintView not implemented")
}
func (UnimplementedPBFTServer) mustEmbedUnimplementedPBFTServer() {}
func (UnimplementedPBFTServer) testEmbeddedByValue()              {}

// UnsafePBFTServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PBFTServer will
// result in compilation errors.
type UnsafePBFTServer interface {
	mustEmbedUnimplementedPBFTServer()
}

func RegisterPBFTServer(s grpc.ServiceRegistrar, srv PBFTServer) {
	// If the following call pancis, it indicates UnimplementedPBFTServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PBFT_ServiceDesc, srv)
}

func _PBFT_Request_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PBFTServer).Request(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PBFT_Request_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PBFTServer).Request(ctx, req.(*RequestMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _PBFT_PrePrepare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrePrepareMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PBFTServer).PrePrepare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PBFT_PrePrepare_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PBFTServer).PrePrepare(ctx, req.(*PrePrepareMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _PBFT_PrePrepared_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AckMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PBFTServer).PrePrepared(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PBFT_PrePrepared_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PBFTServer).PrePrepared(ctx, req.(*AckMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _PBFT_Prepare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AckMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PBFTServer).Prepare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PBFT_Prepare_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PBFTServer).Prepare(ctx, req.(*AckMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _PBFT_Prepared_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AckMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PBFTServer).Prepared(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PBFT_Prepared_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PBFTServer).Prepared(ctx, req.(*AckMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _PBFT_Commit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AckMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PBFTServer).Commit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PBFT_Commit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PBFTServer).Commit(ctx, req.(*AckMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _PBFT_ViewChange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ViewChangeMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PBFTServer).ViewChange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PBFT_ViewChange_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PBFTServer).ViewChange(ctx, req.(*ViewChangeMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _PBFT_NewView_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewViewMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PBFTServer).NewView(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PBFT_NewView_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PBFTServer).NewView(ctx, req.(*NewViewMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _PBFT_Checkpoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckpointMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PBFTServer).Checkpoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PBFT_Checkpoint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PBFTServer).Checkpoint(ctx, req.(*CheckpointMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _PBFT_PrintLog_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PBFTServer).PrintLog(m, &grpc.GenericServerStream[emptypb.Empty, LogRsp]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type PBFT_PrintLogServer = grpc.ServerStreamingServer[LogRsp]

func _PBFT_PrintDB_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PBFTServer).PrintDB(m, &grpc.GenericServerStream[emptypb.Empty, RequestMsg]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type PBFT_PrintDBServer = grpc.ServerStreamingServer[RequestMsg]

func _PBFT_PrintStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PBFTServer).PrintStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PBFT_PrintStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PBFTServer).PrintStatus(ctx, req.(*StatusMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _PBFT_PrintView_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PBFTServer).PrintView(m, &grpc.GenericServerStream[emptypb.Empty, ViewRsp]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type PBFT_PrintViewServer = grpc.ServerStreamingServer[ViewRsp]

// PBFT_ServiceDesc is the grpc.ServiceDesc for PBFT service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PBFT_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pbft.PBFT",
	HandlerType: (*PBFTServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Request",
			Handler:    _PBFT_Request_Handler,
		},
		{
			MethodName: "PrePrepare",
			Handler:    _PBFT_PrePrepare_Handler,
		},
		{
			MethodName: "PrePrepared",
			Handler:    _PBFT_PrePrepared_Handler,
		},
		{
			MethodName: "Prepare",
			Handler:    _PBFT_Prepare_Handler,
		},
		{
			MethodName: "Prepared",
			Handler:    _PBFT_Prepared_Handler,
		},
		{
			MethodName: "Commit",
			Handler:    _PBFT_Commit_Handler,
		},
		{
			MethodName: "ViewChange",
			Handler:    _PBFT_ViewChange_Handler,
		},
		{
			MethodName: "NewView",
			Handler:    _PBFT_NewView_Handler,
		},
		{
			MethodName: "Checkpoint",
			Handler:    _PBFT_Checkpoint_Handler,
		},
		{
			MethodName: "PrintStatus",
			Handler:    _PBFT_PrintStatus_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PrintLog",
			Handler:       _PBFT_PrintLog_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PrintDB",
			Handler:       _PBFT_PrintDB_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PrintView",
			Handler:       _PBFT_PrintView_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pbft.proto",
}
