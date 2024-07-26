// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.2
// source: todo.proto

package todo_protobuf_v1

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

// ToDoServiceClient is the client API for ToDoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ToDoServiceClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponce, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponce, error)
	CreateTask(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*CreateTaskResponce, error)
	ListTasks(ctx context.Context, in *ListTasksRequest, opts ...grpc.CallOption) (*ListTasksResponce, error)
	GetTaskById(ctx context.Context, in *TaskByIdRequest, opts ...grpc.CallOption) (*GetTaskByIdResponce, error)
	UpdateTaskById(ctx context.Context, in *UpdateTaskByIdRequest, opts ...grpc.CallOption) (*ChangedTaskByIdResponce, error)
	DeleteTaskById(ctx context.Context, in *TaskByIdRequest, opts ...grpc.CallOption) (*ChangedTaskByIdResponce, error)
}

type toDoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewToDoServiceClient(cc grpc.ClientConnInterface) ToDoServiceClient {
	return &toDoServiceClient{cc}
}

func (c *toDoServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponce, error) {
	out := new(LoginResponce)
	err := c.cc.Invoke(ctx, "/todo_service.ToDoService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponce, error) {
	out := new(LogoutResponce)
	err := c.cc.Invoke(ctx, "/todo_service.ToDoService/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) CreateTask(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*CreateTaskResponce, error) {
	out := new(CreateTaskResponce)
	err := c.cc.Invoke(ctx, "/todo_service.ToDoService/CreateTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) ListTasks(ctx context.Context, in *ListTasksRequest, opts ...grpc.CallOption) (*ListTasksResponce, error) {
	out := new(ListTasksResponce)
	err := c.cc.Invoke(ctx, "/todo_service.ToDoService/ListTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) GetTaskById(ctx context.Context, in *TaskByIdRequest, opts ...grpc.CallOption) (*GetTaskByIdResponce, error) {
	out := new(GetTaskByIdResponce)
	err := c.cc.Invoke(ctx, "/todo_service.ToDoService/GetTaskById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) UpdateTaskById(ctx context.Context, in *UpdateTaskByIdRequest, opts ...grpc.CallOption) (*ChangedTaskByIdResponce, error) {
	out := new(ChangedTaskByIdResponce)
	err := c.cc.Invoke(ctx, "/todo_service.ToDoService/UpdateTaskById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toDoServiceClient) DeleteTaskById(ctx context.Context, in *TaskByIdRequest, opts ...grpc.CallOption) (*ChangedTaskByIdResponce, error) {
	out := new(ChangedTaskByIdResponce)
	err := c.cc.Invoke(ctx, "/todo_service.ToDoService/DeleteTaskById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ToDoServiceServer is the server API for ToDoService service.
// All implementations must embed UnimplementedToDoServiceServer
// for forward compatibility
type ToDoServiceServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponce, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponce, error)
	CreateTask(context.Context, *CreateTaskRequest) (*CreateTaskResponce, error)
	ListTasks(context.Context, *ListTasksRequest) (*ListTasksResponce, error)
	GetTaskById(context.Context, *TaskByIdRequest) (*GetTaskByIdResponce, error)
	UpdateTaskById(context.Context, *UpdateTaskByIdRequest) (*ChangedTaskByIdResponce, error)
	DeleteTaskById(context.Context, *TaskByIdRequest) (*ChangedTaskByIdResponce, error)
	mustEmbedUnimplementedToDoServiceServer()
}

// UnimplementedToDoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedToDoServiceServer struct {
}

func (UnimplementedToDoServiceServer) Login(context.Context, *LoginRequest) (*LoginResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedToDoServiceServer) Logout(context.Context, *LogoutRequest) (*LogoutResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedToDoServiceServer) CreateTask(context.Context, *CreateTaskRequest) (*CreateTaskResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTask not implemented")
}
func (UnimplementedToDoServiceServer) ListTasks(context.Context, *ListTasksRequest) (*ListTasksResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTasks not implemented")
}
func (UnimplementedToDoServiceServer) GetTaskById(context.Context, *TaskByIdRequest) (*GetTaskByIdResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTaskById not implemented")
}
func (UnimplementedToDoServiceServer) UpdateTaskById(context.Context, *UpdateTaskByIdRequest) (*ChangedTaskByIdResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTaskById not implemented")
}
func (UnimplementedToDoServiceServer) DeleteTaskById(context.Context, *TaskByIdRequest) (*ChangedTaskByIdResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTaskById not implemented")
}
func (UnimplementedToDoServiceServer) mustEmbedUnimplementedToDoServiceServer() {}

// UnsafeToDoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ToDoServiceServer will
// result in compilation errors.
type UnsafeToDoServiceServer interface {
	mustEmbedUnimplementedToDoServiceServer()
}

func RegisterToDoServiceServer(s grpc.ServiceRegistrar, srv ToDoServiceServer) {
	s.RegisterService(&ToDoService_ServiceDesc, srv)
}

func _ToDoService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo_service.ToDoService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo_service.ToDoService/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_CreateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).CreateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo_service.ToDoService/CreateTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).CreateTask(ctx, req.(*CreateTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_ListTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTasksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).ListTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo_service.ToDoService/ListTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).ListTasks(ctx, req.(*ListTasksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_GetTaskById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).GetTaskById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo_service.ToDoService/GetTaskById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).GetTaskById(ctx, req.(*TaskByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_UpdateTaskById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTaskByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).UpdateTaskById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo_service.ToDoService/UpdateTaskById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).UpdateTaskById(ctx, req.(*UpdateTaskByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToDoService_DeleteTaskById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToDoServiceServer).DeleteTaskById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todo_service.ToDoService/DeleteTaskById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToDoServiceServer).DeleteTaskById(ctx, req.(*TaskByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ToDoService_ServiceDesc is the grpc.ServiceDesc for ToDoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ToDoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "todo_service.ToDoService",
	HandlerType: (*ToDoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _ToDoService_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _ToDoService_Logout_Handler,
		},
		{
			MethodName: "CreateTask",
			Handler:    _ToDoService_CreateTask_Handler,
		},
		{
			MethodName: "ListTasks",
			Handler:    _ToDoService_ListTasks_Handler,
		},
		{
			MethodName: "GetTaskById",
			Handler:    _ToDoService_GetTaskById_Handler,
		},
		{
			MethodName: "UpdateTaskById",
			Handler:    _ToDoService_UpdateTaskById_Handler,
		},
		{
			MethodName: "DeleteTaskById",
			Handler:    _ToDoService_DeleteTaskById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "todo.proto",
}