package storage

import (
	"context"
	"errors"

	storageDTO "github.com/IldarGaleev/todo-backend-service/internal/storage/models"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrAccessDenied = errors.New("access denied")
)

// TODO: refactoring is needed - move from here
type IToDoItemProvider interface {
	StorageToDoItem_Create(ctx context.Context, title string, ownerId uint64) (uint64, error)
	StorageToDoItem_DeleteById(ctx context.Context, itemId uint64, ownerId uint64) error
	StorageToDoItem_Update(ctx context.Context, item storageDTO.ToDoItem, ownerId uint64) error
	StorageToDoItem_GetById(ctx context.Context, itemId uint64, ownerId uint64) (*storageDTO.ToDoItem, error)
	StorageToDoItem_GetList(ctx context.Context, ownerId uint64) ([]storageDTO.ToDoItem, error)
}
