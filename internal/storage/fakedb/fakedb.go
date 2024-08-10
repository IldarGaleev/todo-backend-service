// Package fakedb implements fake database provider
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
	dbIndex uint64
}

// Create implements todo.IToDoItemProvider.
func (d *FakeDatabaseProvider) StorageToDoItemCreate(ctx context.Context, title string, ownerID uint64) (uint64, error) {
	d.dbIndex++
	isDone := false
	d.db[d.dbIndex] = storageDTO.ToDoItem{
		Id:         d.dbIndex,
		Title:      &title,
		IsComplete: &isDone,
		OwnerId:    ownerID,
	}

	return d.dbIndex, nil
}

// DeleteById implements todo.IToDoItemProvider.
func (d *FakeDatabaseProvider) StorageToDoItemDeleteByID(ctx context.Context, itemID uint64, ownerID uint64) error {
	if _, ok := d.db[itemID]; ok {
		delete(d.db, itemID)
		return nil
	}
	return storage.ErrNotFound
}

// GetById implements todo.IToDoItemProvider.
func (d *FakeDatabaseProvider) StorageToDoItemGetByID(ctx context.Context, itemID uint64, ownerID uint64) (*storageDTO.ToDoItem, error) {
	if item, ok := d.db[itemID]; ok {
		if item.OwnerId != ownerID {
			return nil, storage.ErrAccessDenied
		}
		return &item, nil
	}
	return nil, storage.ErrNotFound
}

// GetList implements todo.IToDoItemProvider.
func (d *FakeDatabaseProvider) StorageToDoItemGetList(ctx context.Context, ownerID uint64) ([]storageDTO.ToDoItem, error) {
	result := make([]storageDTO.ToDoItem, 0)
	for _, item := range d.db {
		if item.OwnerId == ownerID {
			result = append(result, item)
		}
	}
	return result, nil
}

// Update implements todo.IToDoItemProvider.
func (d *FakeDatabaseProvider) StorageToDoItemUpdate(ctx context.Context, item storageDTO.ToDoItem, ownerID uint64) error {
	if dbItem, ok := d.db[item.Id]; ok {
		if dbItem.OwnerId != ownerID {
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
