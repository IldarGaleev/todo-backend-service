// Package todoservice implements TODO items operations
package todoservice

import (
	"context"
	"errors"
	"log/slog"

	serviceDTO "github.com/IldarGaleev/todo-backend-service/internal/services/servicedto"
	"github.com/IldarGaleev/todo-backend-service/internal/storage"
	storageDTO "github.com/IldarGaleev/todo-backend-service/internal/storage/models"
)

type IToDoItemCreator interface {
	StorageToDoItemCreate(ctx context.Context, title string, ownerID uint64) (uint64, error)
}
type IToDoItemUpdater interface {
	StorageToDoItemUpdate(ctx context.Context, item storageDTO.ToDoItem, ownerID uint64) error
}
type IToDoItemGetter interface {
	StorageToDoItemGetByID(ctx context.Context, itemID uint64, ownerID uint64) (*storageDTO.ToDoItem, error)
	StorageToDoItemGetList(ctx context.Context, ownerID uint64) ([]storageDTO.ToDoItem, error)
}
type IToDoItemDeleter interface {
	StorageToDoItemDeleteByID(ctx context.Context, itemID uint64, ownerID uint64) error
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

func (s *TodoService) Create(ctx context.Context, title string, ownerID uint64) (uint64, error) {
	id, err := s.todoItemsCreator.StorageToDoItemCreate(ctx, title, ownerID)
	if err != nil {
		return 0, errors.Join(ErrInternal, err)
	}
	return id, nil
}

func (s *TodoService) GetByID(ctx context.Context, itemID uint64, ownerID uint64) (*serviceDTO.ToDoItem, error) {
	item, err := s.todoItemsGetter.StorageToDoItemGetByID(ctx, itemID, ownerID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, errors.Join(ErrInternal, err)
	}

	if item.OwnerId != ownerID {
		return nil, ErrAccessDenied
	}

	return &serviceDTO.ToDoItem{
		ID:         itemID,
		OwnerID:    item.OwnerId,
		Title:      item.Title,
		IsComplete: item.IsComplete,
	}, nil
}

func (s *TodoService) GetList(ctx context.Context, ownerID uint64) ([]serviceDTO.ToDoItem, error) {
	storageItems, err := s.todoItemsGetter.StorageToDoItemGetList(ctx, ownerID)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}
	result := make([]serviceDTO.ToDoItem, 0, len(storageItems))

	for _, todoItem := range storageItems {
		result = append(result, serviceDTO.ToDoItem{
			ID:         todoItem.Id,
			OwnerID:    todoItem.OwnerId,
			Title:      todoItem.Title,
			IsComplete: todoItem.IsComplete,
		})
	}

	return result, nil
}

func (s *TodoService) DeleteByID(ctx context.Context, itemID uint64, ownerID uint64) error {
	err := s.todoItemsDeleter.StorageToDoItemDeleteByID(ctx, itemID, ownerID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrItemNotFound
		}
		return errors.Join(ErrInternal, err)
	}

	return nil
}

func (s *TodoService) Update(ctx context.Context, item serviceDTO.ToDoItem, ownerID uint64) error {

	storageItem := storageDTO.ToDoItem{
		Id:         item.ID,
		OwnerId:    item.ID,
		Title:      item.Title,
		IsComplete: item.IsComplete,
	}

	err := s.todoItemsUpdater.StorageToDoItemUpdate(ctx, storageItem, ownerID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrItemNotFound
		}
		return errors.Join(ErrInternal, err)
	}

	return nil
}
