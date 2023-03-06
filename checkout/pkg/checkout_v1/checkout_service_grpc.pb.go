// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: checkout_service.proto

package checkout_v1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CheckoutV1Client is the client API for CheckoutV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CheckoutV1Client interface {
	// Добавить товар в корзину определенного пользователя.
	AddToCart(ctx context.Context, in *AddToCartRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Удалить товар из корзины определенного пользователя.
	DeleteFromCart(ctx context.Context, in *DeleteFromCartRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Показать список товаров в корзине с именами и ценами.
	ListCart(ctx context.Context, in *ListCartRequest, opts ...grpc.CallOption) (*ListCartResponse, error)
	// Оформить заказ по всем товарам корзины.
	Purchase(ctx context.Context, in *PurchaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type checkoutV1Client struct {
	cc grpc.ClientConnInterface
}

func NewCheckoutV1Client(cc grpc.ClientConnInterface) CheckoutV1Client {
	return &checkoutV1Client{cc}
}

func (c *checkoutV1Client) AddToCart(ctx context.Context, in *AddToCartRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/checkout_v1.CheckoutV1/AddToCart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutV1Client) DeleteFromCart(ctx context.Context, in *DeleteFromCartRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/checkout_v1.CheckoutV1/DeleteFromCart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutV1Client) ListCart(ctx context.Context, in *ListCartRequest, opts ...grpc.CallOption) (*ListCartResponse, error) {
	out := new(ListCartResponse)
	err := c.cc.Invoke(ctx, "/checkout_v1.CheckoutV1/ListCart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutV1Client) Purchase(ctx context.Context, in *PurchaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/checkout_v1.CheckoutV1/Purchase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CheckoutV1Server is the server API for CheckoutV1 service.
// All implementations must embed UnimplementedCheckoutV1Server
// for forward compatibility
type CheckoutV1Server interface {
	// Добавить товар в корзину определенного пользователя.
	AddToCart(context.Context, *AddToCartRequest) (*emptypb.Empty, error)
	// Удалить товар из корзины определенного пользователя.
	DeleteFromCart(context.Context, *DeleteFromCartRequest) (*emptypb.Empty, error)
	// Показать список товаров в корзине с именами и ценами.
	ListCart(context.Context, *ListCartRequest) (*ListCartResponse, error)
	// Оформить заказ по всем товарам корзины.
	Purchase(context.Context, *PurchaseRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedCheckoutV1Server()
}

// UnimplementedCheckoutV1Server must be embedded to have forward compatible implementations.
type UnimplementedCheckoutV1Server struct {
}

func (UnimplementedCheckoutV1Server) AddToCart(context.Context, *AddToCartRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddToCart not implemented")
}
func (UnimplementedCheckoutV1Server) DeleteFromCart(context.Context, *DeleteFromCartRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFromCart not implemented")
}
func (UnimplementedCheckoutV1Server) ListCart(context.Context, *ListCartRequest) (*ListCartResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCart not implemented")
}
func (UnimplementedCheckoutV1Server) Purchase(context.Context, *PurchaseRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Purchase not implemented")
}
func (UnimplementedCheckoutV1Server) mustEmbedUnimplementedCheckoutV1Server() {}

// UnsafeCheckoutV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CheckoutV1Server will
// result in compilation errors.
type UnsafeCheckoutV1Server interface {
	mustEmbedUnimplementedCheckoutV1Server()
}

func RegisterCheckoutV1Server(s grpc.ServiceRegistrar, srv CheckoutV1Server) {
	s.RegisterService(&CheckoutV1_ServiceDesc, srv)
}

func _CheckoutV1_AddToCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddToCartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutV1Server).AddToCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/checkout_v1.CheckoutV1/AddToCart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutV1Server).AddToCart(ctx, req.(*AddToCartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckoutV1_DeleteFromCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFromCartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutV1Server).DeleteFromCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/checkout_v1.CheckoutV1/DeleteFromCart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutV1Server).DeleteFromCart(ctx, req.(*DeleteFromCartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckoutV1_ListCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutV1Server).ListCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/checkout_v1.CheckoutV1/ListCart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutV1Server).ListCart(ctx, req.(*ListCartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckoutV1_Purchase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PurchaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutV1Server).Purchase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/checkout_v1.CheckoutV1/Purchase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutV1Server).Purchase(ctx, req.(*PurchaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CheckoutV1_ServiceDesc is the grpc.ServiceDesc for CheckoutV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CheckoutV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "checkout_v1.CheckoutV1",
	HandlerType: (*CheckoutV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddToCart",
			Handler:    _CheckoutV1_AddToCart_Handler,
		},
		{
			MethodName: "DeleteFromCart",
			Handler:    _CheckoutV1_DeleteFromCart_Handler,
		},
		{
			MethodName: "ListCart",
			Handler:    _CheckoutV1_ListCart_Handler,
		},
		{
			MethodName: "Purchase",
			Handler:    _CheckoutV1_Purchase_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "checkout_service.proto",
}
