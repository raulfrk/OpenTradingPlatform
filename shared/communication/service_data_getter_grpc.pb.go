// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: shared/proto/service_data_getter.proto

package communication

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

// DataGetterClient is the client API for DataGetter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataGetterClient interface {
	GetData(ctx context.Context, in *DataRequestExternal, opts ...grpc.CallOption) (*DataResponseExternal, error)
}

type dataGetterClient struct {
	cc grpc.ClientConnInterface
}

func NewDataGetterClient(cc grpc.ClientConnInterface) DataGetterClient {
	return &dataGetterClient{cc}
}

func (c *dataGetterClient) GetData(ctx context.Context, in *DataRequestExternal, opts ...grpc.CallOption) (*DataResponseExternal, error) {
	out := new(DataResponseExternal)
	err := c.cc.Invoke(ctx, "/communication.DataGetter/GetData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataGetterServer is the server API for DataGetter service.
// All implementations must embed UnimplementedDataGetterServer
// for forward compatibility
type DataGetterServer interface {
	GetData(context.Context, *DataRequestExternal) (*DataResponseExternal, error)
	mustEmbedUnimplementedDataGetterServer()
}

// UnimplementedDataGetterServer must be embedded to have forward compatible implementations.
type UnimplementedDataGetterServer struct {
}

func (UnimplementedDataGetterServer) GetData(context.Context, *DataRequestExternal) (*DataResponseExternal, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetData not implemented")
}
func (UnimplementedDataGetterServer) mustEmbedUnimplementedDataGetterServer() {}

// UnsafeDataGetterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataGetterServer will
// result in compilation errors.
type UnsafeDataGetterServer interface {
	mustEmbedUnimplementedDataGetterServer()
}

func RegisterDataGetterServer(s grpc.ServiceRegistrar, srv DataGetterServer) {
	s.RegisterService(&DataGetter_ServiceDesc, srv)
}

func _DataGetter_GetData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataRequestExternal)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataGetterServer).GetData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/communication.DataGetter/GetData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataGetterServer).GetData(ctx, req.(*DataRequestExternal))
	}
	return interceptor(ctx, in, info, handler)
}

// DataGetter_ServiceDesc is the grpc.ServiceDesc for DataGetter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataGetter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "communication.DataGetter",
	HandlerType: (*DataGetterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetData",
			Handler:    _DataGetter_GetData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shared/proto/service_data_getter.proto",
}
