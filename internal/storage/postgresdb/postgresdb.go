package postgresdb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/IldarGaleev/todo-backend-service/internal/storage"
	storageDTO "github.com/IldarGaleev/todo-backend-service/internal/storage/models"
	postgresStorageORM "github.com/IldarGaleev/todo-backend-service/internal/storage/postgresdb/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDataProvider struct {
	log *slog.Logger
	dsn string
	db  *gorm.DB
}

// New create DatabaseApp
func New(log *slog.Logger, dsn string) *PostgresDataProvider {
	return &PostgresDataProvider{
		log: log,
		dsn: dsn,
	}
}

// MustRun create postgres database connection. Panic if failed
func (d *PostgresDataProvider) MustRun() {
	err := d.Run()
	if err != nil {
		panic(err)
	}
}

// Run create postgres database connection
func (d *PostgresDataProvider) Run() error {
	db, err := gorm.Open(postgres.Open(d.dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("database connect error")
	}
	d.db = db
	db.AutoMigrate(
		&postgresStorageORM.UserPG{},
		&postgresStorageORM.ToDoItemPG{},
	)
	return nil
}

// Stop close postgres database connection
func (d *PostgresDataProvider) Stop() error {
	if d.db == nil {
		return storage.ErrDatabaseError
	}

	conn, err := d.db.DB()

	if err != nil {
		return errors.Join(storage.ErrDatabaseError, err)
	}

	err = conn.Close()
	if err != nil {
		return errors.Join(storage.ErrDatabaseError, err)
	}

	return nil
}

// StorageToDoItem_Create implements todoService.IToDoItemCreator.
func (d *PostgresDataProvider) StorageToDoItem_Create(ctx context.Context, title string, ownerId uint64) (uint64, error) {
	newItem := postgresStorageORM.ToDoItemPG{
		OwnerId: ownerId,
		Title:   title,
	}

	result := d.db.WithContext(ctx).Create(&newItem)
	if result.Error != nil {
		return 0, errors.Join(storage.ErrDatabaseError, result.Error)
	}

	return newItem.Id, nil
}

// StorageToDoItem_Update implements todoService.IToDoItemUpdater.
func (d *PostgresDataProvider) StorageToDoItem_Update(ctx context.Context, item storageDTO.ToDoItem, ownerId uint64) error {

	updatedFields := make(map[string]interface{}, 2)

	if item.Title != nil {
		updatedFields["title"] = *item.Title
	}

	if item.IsComplete != nil {
		updatedFields["is_complete"] = *item.IsComplete
	}

	result := d.db.WithContext(ctx).Model(&postgresStorageORM.ToDoItemPG{Id: item.Id}).Updates(updatedFields)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.Join(storage.ErrNotFound, result.Error)
		}
		errors.Join(storage.ErrDatabaseError, result.Error)
	}

	return nil
}

// StorageToDoItem_GetById implements todoService.IToDoItemGetter.
func (d *PostgresDataProvider) StorageToDoItem_GetById(ctx context.Context, itemId uint64, ownerId uint64) (*storageDTO.ToDoItem, error) {
	var item postgresStorageORM.ToDoItemPG
	result := d.db.WithContext(ctx).First(&item, itemId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.Join(storage.ErrNotFound, result.Error)
		}
		return nil, errors.Join(storage.ErrDatabaseError, result.Error)
	}

	return &storageDTO.ToDoItem{
		Id:         item.Id,
		Title:      &item.Title,
		IsComplete: &item.IsComplete,
		OwnerId:    item.OwnerId,
	}, nil
}

// StorageToDoItem_GetList implements todoService.IToDoItemGetter.
func (d *PostgresDataProvider) StorageToDoItem_GetList(ctx context.Context, ownerId uint64) ([]storageDTO.ToDoItem, error) {
	var items []postgresStorageORM.ToDoItemPG
	var resultList []storageDTO.ToDoItem

	result := d.db.WithContext(ctx).Find(&items, "owner_id = ?", ownerId)
	if result.Error != nil {
		return resultList, result.Error
	}

	for _, item := range items {
		resultList = append(resultList, storageDTO.ToDoItem{
			Id:         item.Id,
			Title:      &item.Title,
			IsComplete: &item.IsComplete,
			OwnerId:    item.OwnerId,
		})
	}

	return resultList, nil
}

// StorageToDoItem_DeleteById implements todoService.IToDoItemDeleter.
func (d *PostgresDataProvider) StorageToDoItem_DeleteById(ctx context.Context, itemId uint64, ownerId uint64) error {
	result := d.db.WithContext(ctx).Delete(&postgresStorageORM.ToDoItemPG{}, itemId)

	if result.Error != nil {
		return errors.Join(storage.ErrNotFound, result.Error)
	}

	if result.RowsAffected == 0 {
		return storage.ErrNotFound
	}

	return nil
}

// GetCredential implements credentialService.ICredentialStorageProvider.
func (d *PostgresDataProvider) GetCredential(username string) (*storageDTO.Credential, error) {
	d.log.Warn("get credential not implement")
	return &storageDTO.Credential{
		Username:  "user",
		TokenHash: nil,
	}, nil
}
