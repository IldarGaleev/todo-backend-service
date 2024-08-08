package fakedb

import (
	"context"
	"log/slog"

	"github.com/IldarGaleev/todo-backend-service/internal/storage"
	storageDTO "github.com/IldarGaleev/todo-backend-service/internal/storage/models"
)

type FakeDatabaseProvider struct {
	log      *slog.Logger
	db       map[uint64]storageDTO.ToDoItem
	db_index uint64
}

// Create implements todo.IToDoItemProvider.
func (d *FakeDatabaseProvider) StorageToDoItem_Create(ctx context.Context, title string, ownerId uint64) (uint64, error) {
	d.db_index++
	isDone := false
	d.db[d.db_index] = storageDTO.ToDoItem{
		Id:         d.db_index,
		Title:      &title,
		IsComplete: &isDone,
		OwnerId:    ownerId,
	}

	return d.db_index, nil
}

// DeleteById implements todo.IToDoItemProvider.
func (d *FakeDatabaseProvider) StorageToDoItem_DeleteById(ctx context.Context, itemId uint64, ownerId uint64) error {
	if _, ok := d.db[itemId]; ok {
		delete(d.db, itemId)
		return nil
	}
	return storage.ErrNotFound
}

// GetById implements todo.IToDoItemProvider.
func (d *FakeDatabaseProvider) StorageToDoItem_GetById(ctx context.Context, itemId uint64, ownerId uint64) (*storageDTO.ToDoItem, error) {
	if item, ok := d.db[itemId]; ok {
		if item.OwnerId != ownerId {
			return nil, storage.ErrAccessDenied
		}
		return &item, nil
	}
	return nil, storage.ErrNotFound
}

// GetList implements todo.IToDoItemProvider.
func (d *FakeDatabaseProvider) StorageToDoItem_GetList(ctx context.Context, ownerId uint64) ([]storageDTO.ToDoItem, error) {
	result := make([]storageDTO.ToDoItem, 0)
	for _, item := range d.db {
		if item.OwnerId == ownerId {
			result = append(result, item)
		}
	}
	return result, nil
}

// Update implements todo.IToDoItemProvider.
func (d *FakeDatabaseProvider) StorageToDoItem_Update(ctx context.Context, item storageDTO.ToDoItem, ownerId uint64) error {
	if dbItem, ok := d.db[item.Id]; ok {
		if dbItem.OwnerId != ownerId {
			return storage.ErrAccessDenied
		}

		if item.Title != nil {
			dbItem.Title = item.Title
		}

		if item.IsComplete != nil {
			dbItem.IsComplete = item.IsComplete
		}

		d.db[item.Id] = dbItem

		return nil
	}
	return storage.ErrNotFound
}

// New create DatabaseApp
func New(log *slog.Logger) *FakeDatabaseProvider {
	return &FakeDatabaseProvider{
		log: log.With(slog.String("module","fakedb")),
		db:  make(map[uint64]storageDTO.ToDoItem),
	}
}

// MustRun create postgres database connection. Panic if failed
func (d *FakeDatabaseProvider) MustRun() {
	err := d.Run()
	if err != nil {
		panic(err)
	}
}

// Run create database connection
func (d *FakeDatabaseProvider) Run() error {
	log:=d.log.With(slog.String("method","Run"))
	log.Debug("Start fake DB")
	return nil
}

// Stop close database connection
func (d *FakeDatabaseProvider) Stop() error {
	log:=d.log.With(slog.String("method","Stop"))
	log.Debug("Stop fake DB")
	return nil
}

func (d *FakeDatabaseProvider) GetCredential(username string) (*storageDTO.Credential, error) {
	return &storageDTO.Credential{
		Username:  "user",
		TokenHash: nil,
	}, nil
}

//TODO: unimplement IAccountCreator, IAccountGetter, AuthService
