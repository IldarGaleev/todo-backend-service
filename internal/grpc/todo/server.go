package grpcToDoServer

import (
	"context"

	todo_protobuf_v1 "github.com/IldarGaleev/todo-backend-service/pkg/grpc/proto"
	"google.golang.org/grpc"
)

type serverAPI struct {
	todo_protobuf_v1.UnimplementedToDoServiceServer
}

func Register(gRPC *grpc.Server) {
	todo_protobuf_v1.RegisterToDoServiceServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *todo_protobuf_v1.LoginRequest,
) (*todo_protobuf_v1.LoginResponce, error) {
	return &todo_protobuf_v1.LoginResponce{
		Token: req.GetEmail(),
	}, nil
}

func (s *serverAPI) Logout(
	ctx context.Context,
	req *todo_protobuf_v1.LogoutRequest,
) (*todo_protobuf_v1.LogoutResponce, error) {
	panic("logout not implement")
}

func (s *serverAPI) CreateTask(
	ctx context.Context,
	req *todo_protobuf_v1.CreateTaskRequest,
) (*todo_protobuf_v1.CreateTaskResponce, error) {
	panic("method CreateTask not implemented")
}

func (s *serverAPI) ListTasks(
	ctx context.Context,
	req *todo_protobuf_v1.ListTasksRequest,
) (*todo_protobuf_v1.ListTasksResponce, error) {
	panic("method ListTasks not implemented")
}

func (s *serverAPI) GetTaskById(
	ctx context.Context,
	req *todo_protobuf_v1.TaskByIdRequest,
) (*todo_protobuf_v1.GetTaskByIdResponce, error) {
	panic("method GetTaskById not implemented")
}

func (s *serverAPI) UpdateTaskById(
	ctx context.Context,
	req *todo_protobuf_v1.UpdateTaskByIdRequest,
) (*todo_protobuf_v1.ChangedTaskByIdResponce, error) {
	panic("method UpdateTaskById not implemented")
}

func (s *serverAPI) DeleteTaskById(
	ctx context.Context,
	req *todo_protobuf_v1.TaskByIdRequest,
) (*todo_protobuf_v1.ChangedTaskByIdResponce, error) {
	panic("method DeleteTaskById not implemented")
}
