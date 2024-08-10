// TODO: split to auth and todo items grpc....

// Package grpctodoserver implements gRPC handlers
package grpctodoserver

import (
	"context"

	serviceDTO "github.com/IldarGaleev/todo-backend-service/internal/services/servicedto"
	todo_protobuf_v1 "github.com/IldarGaleev/todo-backend-service/pkg/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IToDoItemCreatorService interface {
	Create(ctx context.Context, title string, ownerID uint64) (uint64, error)
}

type IToDoItemGetterService interface {
	GetByID(ctx context.Context, itemID uint64, ownerID uint64) (*serviceDTO.ToDoItem, error)
	GetList(ctx context.Context, ownerID uint64) ([]serviceDTO.ToDoItem, error)
}

type IToDoItemDeleterService interface {
	DeleteByID(ctx context.Context, itemID uint64, ownerID uint64) error
}

type IToDoItemUpdaterService interface {
	Update(ctx context.Context, item serviceDTO.ToDoItem, ownerID uint64) error
}

type IAccountSecretCreator interface {
	CreateUserSecret(ctx context.Context, user serviceDTO.User) (string, error)
}

type IAccountSecretValidator interface {
	CheckSecret(ctx context.Context, secret []byte) (
		*serviceDTO.User,
		error,
	)
}
type IAccountSecretDeleter interface {
	DeleteSecret(ctx context.Context, secret []byte) error
}
type serverAPI struct {
	todo_protobuf_v1.UnimplementedToDoServiceServer
	todoItemsCreatorService IToDoItemCreatorService
	todoItemsUpdaterService IToDoItemUpdaterService
	todoItemsGetterService  IToDoItemGetterService
	todoItemsDeleterService IToDoItemDeleterService
	accountSecretCreator    IAccountSecretCreator
	accountSecretValidator  IAccountSecretValidator
	accountSecretDeleter    IAccountSecretDeleter
}

func Register(
	gRPC *grpc.Server,
	todoItemsCreatorService IToDoItemCreatorService,
	todoItemsUpdaterService IToDoItemUpdaterService,
	todoItemsGetterService IToDoItemGetterService,
	todoItemsDeleterService IToDoItemDeleterService,
	accountSecretCreator IAccountSecretCreator,
	accountSecretValidator IAccountSecretValidator,
	accountSecretDeleter IAccountSecretDeleter,
) {
	todo_protobuf_v1.RegisterToDoServiceServer(
		gRPC,
		&serverAPI{
			todoItemsCreatorService: todoItemsCreatorService,
			todoItemsUpdaterService: todoItemsUpdaterService,
			todoItemsGetterService:  todoItemsGetterService,
			todoItemsDeleterService: todoItemsDeleterService,
			accountSecretCreator:    accountSecretCreator,
			accountSecretValidator:  accountSecretValidator,
			accountSecretDeleter:    accountSecretDeleter,
		},
	)
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *todo_protobuf_v1.LoginRequest,
) (*todo_protobuf_v1.LoginResponce, error) {

	token, err := s.accountSecretCreator.CreateUserSecret(
		ctx,
		serviceDTO.User{
			Username: &req.Email,
			Password: req.Password,
		},
	)

	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "wrong username or password")
	}

	return &todo_protobuf_v1.LoginResponce{
		Token: token,
	}, nil
}

func (s *serverAPI) Logout(
	ctx context.Context,
	req *todo_protobuf_v1.LogoutRequest,
) (*todo_protobuf_v1.LogoutResponce, error) {
	t := req.GetToken()
	tt := req.Token
	_ = tt
	err := s.accountSecretDeleter.DeleteSecret(ctx, []byte(t))
	if err != nil {
		return &todo_protobuf_v1.LogoutResponce{
			Success: false,
		}, status.Error(codes.FailedPrecondition, "wrong token or revoked")
	}

	return &todo_protobuf_v1.LogoutResponce{
		Success: true,
	}, nil
}

func (s *serverAPI) CheckSecret(
	ctx context.Context,
	req *todo_protobuf_v1.CheckSecretRequest,
) (*todo_protobuf_v1.CheckSecretResponce, error) {
	user, err := s.accountSecretValidator.CheckSecret(
		ctx,
		[]byte(req.GetSecret()),
	)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "check secret failed")
	}

	return &todo_protobuf_v1.CheckSecretResponce{
		UserId: *user.UserID,
		Email:  *user.Username,
	}, nil
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
			TaskId: item.ID,
			Title:  *item.Title,
			IsDone: *item.IsComplete,
		})
	}
	return &todo_protobuf_v1.ListTasksResponce{
		Tasks: responseItems,
	}, nil
}

func (s *serverAPI) GetTaskByID(
	ctx context.Context,
	req *todo_protobuf_v1.TaskByIdRequest,
) (*todo_protobuf_v1.GetTaskByIdResponce, error) {
	item, err := s.todoItemsGetterService.GetByID(ctx, req.GetTaskId(), req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "Item not found")
	}

	return &todo_protobuf_v1.GetTaskByIdResponce{
		TaskId: req.TaskId,
		Title:  *item.Title,
		IsDone: *item.IsComplete,
	}, nil
}

func (s *serverAPI) UpdateTaskByID(
	ctx context.Context,
	req *todo_protobuf_v1.UpdateTaskByIdRequest,
) (*todo_protobuf_v1.ChangedTaskByIdResponce, error) {

	err := s.todoItemsUpdaterService.Update(ctx, serviceDTO.ToDoItem{
		ID:         req.GetTaskId(),
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

func (s *serverAPI) DeleteTaskByID(
	ctx context.Context,
	req *todo_protobuf_v1.TaskByIdRequest,
) (*todo_protobuf_v1.ChangedTaskByIdResponce, error) {
	err := s.todoItemsDeleterService.DeleteByID(ctx, req.GetTaskId(), req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "Item not found")
	}

	return &todo_protobuf_v1.ChangedTaskByIdResponce{
		TaskId:    req.GetTaskId(),
		IsSuccess: true,
	}, nil
}
