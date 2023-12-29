// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: word_service.proto

package pb

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

// WordServiceClient is the client API for WordService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WordServiceClient interface {
	FindWordByDictionary(ctx context.Context, in *FindWordByDictionaryRequest, opts ...grpc.CallOption) (*FindWordByDictionaryResponse, error)
	CreateFavoriteWordMeaning(ctx context.Context, in *CreateFavoriteWordMeaningRequest, opts ...grpc.CallOption) (*CreateFavoriteWordMeaningResponse, error)
}

type wordServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWordServiceClient(cc grpc.ClientConnInterface) WordServiceClient {
	return &wordServiceClient{cc}
}

func (c *wordServiceClient) FindWordByDictionary(ctx context.Context, in *FindWordByDictionaryRequest, opts ...grpc.CallOption) (*FindWordByDictionaryResponse, error) {
	out := new(FindWordByDictionaryResponse)
	err := c.cc.Invoke(ctx, "/pb.WordService/FindWordByDictionary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wordServiceClient) CreateFavoriteWordMeaning(ctx context.Context, in *CreateFavoriteWordMeaningRequest, opts ...grpc.CallOption) (*CreateFavoriteWordMeaningResponse, error) {
	out := new(CreateFavoriteWordMeaningResponse)
	err := c.cc.Invoke(ctx, "/pb.WordService/CreateFavoriteWordMeaning", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WordServiceServer is the server API for WordService service.
// All implementations must embed UnimplementedWordServiceServer
// for forward compatibility
type WordServiceServer interface {
	FindWordByDictionary(context.Context, *FindWordByDictionaryRequest) (*FindWordByDictionaryResponse, error)
	CreateFavoriteWordMeaning(context.Context, *CreateFavoriteWordMeaningRequest) (*CreateFavoriteWordMeaningResponse, error)
	mustEmbedUnimplementedWordServiceServer()
}

// UnimplementedWordServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWordServiceServer struct {
}

func (UnimplementedWordServiceServer) FindWordByDictionary(context.Context, *FindWordByDictionaryRequest) (*FindWordByDictionaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindWordByDictionary not implemented")
}
func (UnimplementedWordServiceServer) CreateFavoriteWordMeaning(context.Context, *CreateFavoriteWordMeaningRequest) (*CreateFavoriteWordMeaningResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFavoriteWordMeaning not implemented")
}
func (UnimplementedWordServiceServer) mustEmbedUnimplementedWordServiceServer() {}

// UnsafeWordServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WordServiceServer will
// result in compilation errors.
type UnsafeWordServiceServer interface {
	mustEmbedUnimplementedWordServiceServer()
}

func RegisterWordServiceServer(s grpc.ServiceRegistrar, srv WordServiceServer) {
	s.RegisterService(&WordService_ServiceDesc, srv)
}

func _WordService_FindWordByDictionary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindWordByDictionaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WordServiceServer).FindWordByDictionary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.WordService/FindWordByDictionary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WordServiceServer).FindWordByDictionary(ctx, req.(*FindWordByDictionaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WordService_CreateFavoriteWordMeaning_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFavoriteWordMeaningRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WordServiceServer).CreateFavoriteWordMeaning(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.WordService/CreateFavoriteWordMeaning",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WordServiceServer).CreateFavoriteWordMeaning(ctx, req.(*CreateFavoriteWordMeaningRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WordService_ServiceDesc is the grpc.ServiceDesc for WordService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WordService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.WordService",
	HandlerType: (*WordServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindWordByDictionary",
			Handler:    _WordService_FindWordByDictionary_Handler,
		},
		{
			MethodName: "CreateFavoriteWordMeaning",
			Handler:    _WordService_CreateFavoriteWordMeaning_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "word_service.proto",
}
