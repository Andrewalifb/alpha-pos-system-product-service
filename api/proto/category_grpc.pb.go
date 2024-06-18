// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: category.proto

package lpha_pos_system_product_service

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

// PosProductCategoryServiceClient is the client API for PosProductCategoryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PosProductCategoryServiceClient interface {
	CreatePosProductCategory(ctx context.Context, in *CreatePosProductCategoryRequest, opts ...grpc.CallOption) (*CreatePosProductCategoryResponse, error)
	ReadPosProductCategory(ctx context.Context, in *ReadPosProductCategoryRequest, opts ...grpc.CallOption) (*ReadPosProductCategoryResponse, error)
	UpdatePosProductCategory(ctx context.Context, in *UpdatePosProductCategoryRequest, opts ...grpc.CallOption) (*UpdatePosProductCategoryResponse, error)
	DeletePosProductCategory(ctx context.Context, in *DeletePosProductCategoryRequest, opts ...grpc.CallOption) (*DeletePosProductCategoryResponse, error)
	ReadAllPosProductCategories(ctx context.Context, in *ReadAllPosProductCategoriesRequest, opts ...grpc.CallOption) (*ReadAllPosProductCategoriesResponse, error)
}

type posProductCategoryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPosProductCategoryServiceClient(cc grpc.ClientConnInterface) PosProductCategoryServiceClient {
	return &posProductCategoryServiceClient{cc}
}

func (c *posProductCategoryServiceClient) CreatePosProductCategory(ctx context.Context, in *CreatePosProductCategoryRequest, opts ...grpc.CallOption) (*CreatePosProductCategoryResponse, error) {
	out := new(CreatePosProductCategoryResponse)
	err := c.cc.Invoke(ctx, "/pos.PosProductCategoryService/CreatePosProductCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *posProductCategoryServiceClient) ReadPosProductCategory(ctx context.Context, in *ReadPosProductCategoryRequest, opts ...grpc.CallOption) (*ReadPosProductCategoryResponse, error) {
	out := new(ReadPosProductCategoryResponse)
	err := c.cc.Invoke(ctx, "/pos.PosProductCategoryService/ReadPosProductCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *posProductCategoryServiceClient) UpdatePosProductCategory(ctx context.Context, in *UpdatePosProductCategoryRequest, opts ...grpc.CallOption) (*UpdatePosProductCategoryResponse, error) {
	out := new(UpdatePosProductCategoryResponse)
	err := c.cc.Invoke(ctx, "/pos.PosProductCategoryService/UpdatePosProductCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *posProductCategoryServiceClient) DeletePosProductCategory(ctx context.Context, in *DeletePosProductCategoryRequest, opts ...grpc.CallOption) (*DeletePosProductCategoryResponse, error) {
	out := new(DeletePosProductCategoryResponse)
	err := c.cc.Invoke(ctx, "/pos.PosProductCategoryService/DeletePosProductCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *posProductCategoryServiceClient) ReadAllPosProductCategories(ctx context.Context, in *ReadAllPosProductCategoriesRequest, opts ...grpc.CallOption) (*ReadAllPosProductCategoriesResponse, error) {
	out := new(ReadAllPosProductCategoriesResponse)
	err := c.cc.Invoke(ctx, "/pos.PosProductCategoryService/ReadAllPosProductCategories", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PosProductCategoryServiceServer is the server API for PosProductCategoryService service.
// All implementations must embed UnimplementedPosProductCategoryServiceServer
// for forward compatibility
type PosProductCategoryServiceServer interface {
	CreatePosProductCategory(context.Context, *CreatePosProductCategoryRequest) (*CreatePosProductCategoryResponse, error)
	ReadPosProductCategory(context.Context, *ReadPosProductCategoryRequest) (*ReadPosProductCategoryResponse, error)
	UpdatePosProductCategory(context.Context, *UpdatePosProductCategoryRequest) (*UpdatePosProductCategoryResponse, error)
	DeletePosProductCategory(context.Context, *DeletePosProductCategoryRequest) (*DeletePosProductCategoryResponse, error)
	ReadAllPosProductCategories(context.Context, *ReadAllPosProductCategoriesRequest) (*ReadAllPosProductCategoriesResponse, error)
	mustEmbedUnimplementedPosProductCategoryServiceServer()
}

// UnimplementedPosProductCategoryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPosProductCategoryServiceServer struct {
}

func (UnimplementedPosProductCategoryServiceServer) CreatePosProductCategory(context.Context, *CreatePosProductCategoryRequest) (*CreatePosProductCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePosProductCategory not implemented")
}
func (UnimplementedPosProductCategoryServiceServer) ReadPosProductCategory(context.Context, *ReadPosProductCategoryRequest) (*ReadPosProductCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadPosProductCategory not implemented")
}
func (UnimplementedPosProductCategoryServiceServer) UpdatePosProductCategory(context.Context, *UpdatePosProductCategoryRequest) (*UpdatePosProductCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePosProductCategory not implemented")
}
func (UnimplementedPosProductCategoryServiceServer) DeletePosProductCategory(context.Context, *DeletePosProductCategoryRequest) (*DeletePosProductCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePosProductCategory not implemented")
}
func (UnimplementedPosProductCategoryServiceServer) ReadAllPosProductCategories(context.Context, *ReadAllPosProductCategoriesRequest) (*ReadAllPosProductCategoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadAllPosProductCategories not implemented")
}
func (UnimplementedPosProductCategoryServiceServer) mustEmbedUnimplementedPosProductCategoryServiceServer() {
}

// UnsafePosProductCategoryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PosProductCategoryServiceServer will
// result in compilation errors.
type UnsafePosProductCategoryServiceServer interface {
	mustEmbedUnimplementedPosProductCategoryServiceServer()
}

func RegisterPosProductCategoryServiceServer(s grpc.ServiceRegistrar, srv PosProductCategoryServiceServer) {
	s.RegisterService(&PosProductCategoryService_ServiceDesc, srv)
}

func _PosProductCategoryService_CreatePosProductCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePosProductCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PosProductCategoryServiceServer).CreatePosProductCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pos.PosProductCategoryService/CreatePosProductCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PosProductCategoryServiceServer).CreatePosProductCategory(ctx, req.(*CreatePosProductCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PosProductCategoryService_ReadPosProductCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadPosProductCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PosProductCategoryServiceServer).ReadPosProductCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pos.PosProductCategoryService/ReadPosProductCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PosProductCategoryServiceServer).ReadPosProductCategory(ctx, req.(*ReadPosProductCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PosProductCategoryService_UpdatePosProductCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePosProductCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PosProductCategoryServiceServer).UpdatePosProductCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pos.PosProductCategoryService/UpdatePosProductCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PosProductCategoryServiceServer).UpdatePosProductCategory(ctx, req.(*UpdatePosProductCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PosProductCategoryService_DeletePosProductCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePosProductCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PosProductCategoryServiceServer).DeletePosProductCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pos.PosProductCategoryService/DeletePosProductCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PosProductCategoryServiceServer).DeletePosProductCategory(ctx, req.(*DeletePosProductCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PosProductCategoryService_ReadAllPosProductCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadAllPosProductCategoriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PosProductCategoryServiceServer).ReadAllPosProductCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pos.PosProductCategoryService/ReadAllPosProductCategories",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PosProductCategoryServiceServer).ReadAllPosProductCategories(ctx, req.(*ReadAllPosProductCategoriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PosProductCategoryService_ServiceDesc is the grpc.ServiceDesc for PosProductCategoryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PosProductCategoryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pos.PosProductCategoryService",
	HandlerType: (*PosProductCategoryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePosProductCategory",
			Handler:    _PosProductCategoryService_CreatePosProductCategory_Handler,
		},
		{
			MethodName: "ReadPosProductCategory",
			Handler:    _PosProductCategoryService_ReadPosProductCategory_Handler,
		},
		{
			MethodName: "UpdatePosProductCategory",
			Handler:    _PosProductCategoryService_UpdatePosProductCategory_Handler,
		},
		{
			MethodName: "DeletePosProductCategory",
			Handler:    _PosProductCategoryService_DeletePosProductCategory_Handler,
		},
		{
			MethodName: "ReadAllPosProductCategories",
			Handler:    _PosProductCategoryService_ReadAllPosProductCategories_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "category.proto",
}
