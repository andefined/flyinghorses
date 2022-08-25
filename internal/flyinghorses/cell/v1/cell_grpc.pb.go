// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package cellv1

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

// CellServiceClient is the client API for CellService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CellServiceClient interface {
	Stream(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (CellService_StreamClient, error)
}

type cellServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCellServiceClient(cc grpc.ClientConnInterface) CellServiceClient {
	return &cellServiceClient{cc}
}

func (c *cellServiceClient) Stream(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (CellService_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &CellService_ServiceDesc.Streams[0], "/flyinghorses.cell.v1.CellService/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &cellServiceStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CellService_StreamClient interface {
	Recv() (*StreamResponse, error)
	grpc.ClientStream
}

type cellServiceStreamClient struct {
	grpc.ClientStream
}

func (x *cellServiceStreamClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CellServiceServer is the server API for CellService service.
// All implementations should embed UnimplementedCellServiceServer
// for forward compatibility
type CellServiceServer interface {
	Stream(*StreamRequest, CellService_StreamServer) error
}

// UnimplementedCellServiceServer should be embedded to have forward compatible implementations.
type UnimplementedCellServiceServer struct {
}

func (UnimplementedCellServiceServer) Stream(*StreamRequest, CellService_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}

// UnsafeCellServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CellServiceServer will
// result in compilation errors.
type UnsafeCellServiceServer interface {
	mustEmbedUnimplementedCellServiceServer()
}

func RegisterCellServiceServer(s grpc.ServiceRegistrar, srv CellServiceServer) {
	s.RegisterService(&CellService_ServiceDesc, srv)
}

func _CellService_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CellServiceServer).Stream(m, &cellServiceStreamServer{stream})
}

type CellService_StreamServer interface {
	Send(*StreamResponse) error
	grpc.ServerStream
}

type cellServiceStreamServer struct {
	grpc.ServerStream
}

func (x *cellServiceStreamServer) Send(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// CellService_ServiceDesc is the grpc.ServiceDesc for CellService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CellService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "flyinghorses.cell.v1.CellService",
	HandlerType: (*CellServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _CellService_Stream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "flyinghorses/cell/v1/cell.proto",
}