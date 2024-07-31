package todoService

import (
	"context"
	"log/slog"

	"github.com/IldarGaleev/todo-backend-service/internal/domain/models"
	"github.com/IldarGaleev/todo-backend-service/internal/storage"
)

type TodoService struct {
	logger           *slog.Logger
	todoItemsStorage storage.IToDoItemProvider
}

func New(log *slog.Logger, todoItemsStorageProvider storage.IToDoItemProvider) *TodoService {
	return &TodoService{
		logger:           log,
		todoItemsStorage: todoItemsStorageProvider,
	}
}

func (s *TodoService) Create(ctx context.Context, title string, ownerId uint64) (uint64, error) {
	id, err := s.todoItemsStorage.Create(ctx, title, ownerId)
	if err != nil {
		//TODO: wrap error
		return 0, err
	}
	return id, nil
}

func (s *TodoService) GetById(ctx context.Context, itemId uint64, ownerId uint64) (*models.ToDoItem, error) {
	item, err := s.todoItemsStorage.GetById(ctx, itemId, ownerId)
	if err != nil {
		//TODO: wrap error
		return nil, err
	}
	return item, nil
}

func (s *TodoService) GetList(ctx context.Context, ownerId uint64) ([]models.ToDoItem, error) {
	items, err := s.todoItemsStorage.GetList(ctx, ownerId)
	if err != nil {
		//TODO: wrap error
		return nil, err
	}
	return items, nil
}

func (s *TodoService) DeleteById(ctx context.Context, itemId uint64, ownerId uint64) error {
	err := s.todoItemsStorage.DeleteById(ctx, itemId, ownerId)
	if err != nil {
		//TODO: wrap error
		return err
	}

	return nil
}

func (s *TodoService) Update(ctx context.Context, item models.ToDoItem, ownerId uint64) error {
	err := s.todoItemsStorage.Update(ctx, item, ownerId)
	if err != nil {
		//TODO: wrap error
		return err
	}

	return nil
}
