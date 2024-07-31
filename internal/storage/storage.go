package storage

import (
	"context"
	"errors"

	"github.com/IldarGaleev/todo-backend-service/internal/domain/models"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrAccessDenied = errors.New("access denied")
)

// TODO: refactoring is needed - move from here
type IToDoItemProvider interface {
	Create(ctx context.Context, title string, ownerId uint64) (uint64, error)
	DeleteById(ctx context.Context, itemId uint64, ownerId uint64) error
	Update(ctx context.Context, item models.ToDoItem, ownerId uint64) error
	GetById(ctx context.Context, itemId uint64, ownerId uint64) (*models.ToDoItem, error)
	GetList(ctx context.Context, ownerId uint64) ([]models.ToDoItem, error)
}
