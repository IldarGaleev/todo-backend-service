// TODO: split to auth and todo items grpc....
package grpcToDoServer

import (
	"context"

	serviceDTO "github.com/IldarGaleev/todo-backend-service/internal/services/models"
	todo_protobuf_v1 "github.com/IldarGaleev/todo-backend-service/pkg/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IToDoItemCreatorService interface {
	Create(ctx context.Context, title string, ownerId uint64) (uint64, error)
}

type IToDoItemGetterService interface {
	GetById(ctx context.Context, itemId uint64, ownerId uint64) (*serviceDTO.ToDoItem, error)
	GetList(ctx context.Context, ownerId uint64) ([]serviceDTO.ToDoItem, error)
}

type IToDoItemDeleterService interface {
	DeleteById(ctx context.Context, itemId uint64, ownerId uint64) error
}

type IToDoItemUpdaterService interface {
	Update(ctx context.Context, item serviceDTO.ToDoItem, ownerId uint64) error
}

type serverAPI struct {
	todo_protobuf_v1.UnimplementedToDoServiceServer
	todoItemsCreatorService IToDoItemCreatorService
	todoItemsUpdaterService IToDoItemUpdaterService
	todoItemsGetterService  IToDoItemGetterService
	todoItemsDeleterService IToDoItemDeleterService
}

func Register(
	gRPC *grpc.Server,
	todoItemsCreatorService IToDoItemCreatorService,
	todoItemsUpdaterService IToDoItemUpdaterService,
	todoItemsGetterService IToDoItemGetterService,
	todoItemsDeleterService IToDoItemDeleterService,
) {
	todo_protobuf_v1.RegisterToDoServiceServer(
		gRPC,
		&serverAPI{
			todoItemsCreatorService: todoItemsCreatorService,
			todoItemsUpdaterService: todoItemsUpdaterService,
			todoItemsGetterService:  todoItemsGetterService,
			todoItemsDeleterService: todoItemsDeleterService,
		},
	)
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *todo_protobuf_v1.LoginRequest,
) (*todo_protobuf_v1.LoginResponce, error) {
	panic("login not implement")
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
	id, err := s.todoItemsCreatorService.Create(ctx, req.GetTitle(), req.GetUserId())

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal create error")
	}

	return &todo_protobuf_v1.CreateTaskResponce{
		TaskId: id,
	}, nil

}

func (s *serverAPI) ListTasks(
	ctx context.Context,
	req *todo_protobuf_v1.ListTasksRequest,
) (*todo_protobuf_v1.ListTasksResponce, error) {
	items, err := s.todoItemsGetterService.GetList(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	responseItems := make([]*todo_protobuf_v1.GetTaskByIdResponce, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, &todo_protobuf_v1.GetTaskByIdResponce{
			TaskId: item.Id,
			Title:  *item.Title,
			IsDone: *item.IsComplete,
		})
	}
	return &todo_protobuf_v1.ListTasksResponce{
		Tasks: responseItems,
	}, nil
}

func (s *serverAPI) GetTaskById(
	ctx context.Context,
	req *todo_protobuf_v1.TaskByIdRequest,
) (*todo_protobuf_v1.GetTaskByIdResponce, error) {
	item, err := s.todoItemsGetterService.GetById(ctx, req.GetTaskId(), req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "Item not found")
	}

	return &todo_protobuf_v1.GetTaskByIdResponce{
		TaskId: req.TaskId,
		Title:  *item.Title,
		IsDone: *item.IsComplete,
	}, nil
}

func (s *serverAPI) UpdateTaskById(
	ctx context.Context,
	req *todo_protobuf_v1.UpdateTaskByIdRequest,
) (*todo_protobuf_v1.ChangedTaskByIdResponce, error) {

	err := s.todoItemsUpdaterService.Update(ctx, serviceDTO.ToDoItem{
		Id:         req.GetTaskId(),
		Title:      req.Title,
		IsComplete: req.IsDone,
	}, req.GetUserId())

	if err != nil {
		return nil, status.Error(codes.NotFound, "Item not found")
	}

	return &todo_protobuf_v1.ChangedTaskByIdResponce{
		TaskId:    req.GetTaskId(),
		IsSuccess: true,
	}, nil
}

func (s *serverAPI) DeleteTaskById(
	ctx context.Context,
	req *todo_protobuf_v1.TaskByIdRequest,
) (*todo_protobuf_v1.ChangedTaskByIdResponce, error) {
	err := s.todoItemsDeleterService.DeleteById(ctx, req.GetTaskId(), req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "Item not found")
	}

	return &todo_protobuf_v1.ChangedTaskByIdResponce{
		TaskId:    req.GetTaskId(),
		IsSuccess: true,
	}, nil
}
