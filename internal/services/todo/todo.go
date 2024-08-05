package todoService

import (
	"context"
	"log/slog"

	serviceDTO "github.com/IldarGaleev/todo-backend-service/internal/services/models"
	storageDTO "github.com/IldarGaleev/todo-backend-service/internal/storage/models"
)

type IToDoItemCreator interface {
	StorageToDoItem_Create(ctx context.Context, title string, ownerId uint64) (uint64, error)
}
type IToDoItemUpdater interface {
	StorageToDoItem_Update(ctx context.Context, item storageDTO.ToDoItem, ownerId uint64) error
}
type IToDoItemGetter interface {
	StorageToDoItem_GetById(ctx context.Context, itemId uint64, ownerId uint64) (*storageDTO.ToDoItem, error)
	StorageToDoItem_GetList(ctx context.Context, ownerId uint64) ([]storageDTO.ToDoItem, error)
}
type IToDoItemDeleter interface {
	StorageToDoItem_DeleteById(ctx context.Context, itemId uint64, ownerId uint64) error
}

type TodoService struct {
	logger           *slog.Logger
	todoItemsCreator IToDoItemCreator
	todoItemsUpdater IToDoItemUpdater
	todoItemsGetter  IToDoItemGetter
	todoItemsDeleter IToDoItemDeleter
}

func New(
	log *slog.Logger,
	todoItemsCreator IToDoItemCreator,
	todoItemsUpdater IToDoItemUpdater,
	todoItemsGetter IToDoItemGetter,
	todoItemsDeleter IToDoItemDeleter,
) *TodoService {
	return &TodoService{
		logger:           log,
		todoItemsCreator: todoItemsCreator,
		todoItemsUpdater: todoItemsUpdater,
		todoItemsGetter:  todoItemsGetter,
		todoItemsDeleter: todoItemsDeleter,
	}
}

func (s *TodoService) Create(ctx context.Context, title string, ownerId uint64) (uint64, error) {
	id, err := s.todoItemsCreator.StorageToDoItem_Create(ctx, title, ownerId)
	if err != nil {
		//TODO: wrap error
		return 0, err
	}
	return id, nil
}

func (s *TodoService) GetById(ctx context.Context, itemId uint64, ownerId uint64) (*serviceDTO.ToDoItem, error) {
	item, err := s.todoItemsGetter.StorageToDoItem_GetById(ctx, itemId, ownerId)
	if err != nil {
		//TODO: wrap error
		return nil, err
	}
	return &serviceDTO.ToDoItem{
		Id:         itemId,
		OwnerId:    item.OwnerId,
		Title:      item.Title,
		IsComplete: item.IsComplete,
	}, nil
}

func (s *TodoService) GetList(ctx context.Context, ownerId uint64) ([]serviceDTO.ToDoItem, error) {
	storageItems, err := s.todoItemsGetter.StorageToDoItem_GetList(ctx, ownerId)
	if err != nil {
		//TODO: wrap error
		return nil, err
	}
	result := make([]serviceDTO.ToDoItem, 0, len(storageItems))

	for _, todoItem := range storageItems {
		result = append(result, serviceDTO.ToDoItem{
			Id:         todoItem.Id,
			OwnerId:    todoItem.OwnerId,
			Title:      todoItem.Title,
			IsComplete: todoItem.IsComplete,
		})
	}

	return result, nil
}

func (s *TodoService) DeleteById(ctx context.Context, itemId uint64, ownerId uint64) error {
	err := s.todoItemsDeleter.StorageToDoItem_DeleteById(ctx, itemId, ownerId)
	if err != nil {
		//TODO: wrap error
		return err
	}

	return nil
}

func (s *TodoService) Update(ctx context.Context, item serviceDTO.ToDoItem, ownerId uint64) error {

	storageItem := storageDTO.ToDoItem{
		Id:         item.Id,
		OwnerId:    item.Id,
		Title:      item.Title,
		IsComplete: item.IsComplete,
	}

	err := s.todoItemsUpdater.StorageToDoItem_Update(ctx, storageItem, ownerId)
	if err != nil {
		//TODO: wrap error
		return err
	}

	return nil
}
