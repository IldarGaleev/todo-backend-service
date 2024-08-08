package todoService

import (
	"context"
	"errors"
	"log/slog"

	serviceDTO "github.com/IldarGaleev/todo-backend-service/internal/services/models"
	"github.com/IldarGaleev/todo-backend-service/internal/storage"
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

var (
	ErrAccessDenied = errors.New("todo service: access denied")
	ErrItemNotFound = errors.New("todo service: item not found")
	ErrInternal     = errors.New("todo service: internal error")
)

func New(
	log *slog.Logger,
	todoItemsCreator IToDoItemCreator,
	todoItemsUpdater IToDoItemUpdater,
	todoItemsGetter IToDoItemGetter,
	todoItemsDeleter IToDoItemDeleter,
) *TodoService {
	return &TodoService{
		logger:           log.With(slog.String("module", "todoService")),
		todoItemsCreator: todoItemsCreator,
		todoItemsUpdater: todoItemsUpdater,
		todoItemsGetter:  todoItemsGetter,
		todoItemsDeleter: todoItemsDeleter,
	}
}

func (s *TodoService) Create(ctx context.Context, title string, ownerId uint64) (uint64, error) {
	id, err := s.todoItemsCreator.StorageToDoItem_Create(ctx, title, ownerId)
	if err != nil {
		return 0, errors.Join(ErrInternal, err)
	}
	return id, nil
}

func (s *TodoService) GetById(ctx context.Context, itemId uint64, ownerId uint64) (*serviceDTO.ToDoItem, error) {
	item, err := s.todoItemsGetter.StorageToDoItem_GetById(ctx, itemId, ownerId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, errors.Join(ErrInternal, err)
	}

	if item.OwnerId != ownerId {
		return nil, ErrAccessDenied
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
		return nil, errors.Join(ErrInternal, err)
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
		if errors.Is(err, storage.ErrNotFound) {
			return ErrItemNotFound
		}
		return errors.Join(ErrInternal, err)
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
		if errors.Is(err, storage.ErrNotFound) {
			return ErrItemNotFound
		}
		return errors.Join(ErrInternal, err)
	}

	return nil
}
