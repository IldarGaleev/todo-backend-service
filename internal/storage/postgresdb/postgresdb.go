// Package postgresdb implements Postgres data provider
package postgresdb

import (
	"context"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
	"log/slog"

	serviceDTO "github.com/IldarGaleev/todo-backend-service/internal/services/servicedto"
	"github.com/IldarGaleev/todo-backend-service/internal/storage"
	storageDTO "github.com/IldarGaleev/todo-backend-service/internal/storage/models"
	postgresStorageORM "github.com/IldarGaleev/todo-backend-service/internal/storage/postgresdb/postgresstorageorm"
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
		log: log.With(slog.String("module", "postgresdb")),
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
func (d *PostgresDataProvider) runWithDialector(dialector gorm.Dialector, silentLog bool) error {
	db, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if silentLog {
		db.Config.Logger = logger.Default.LogMode(logger.Silent)
	}

	if err != nil {
		return errors.Join(storage.ErrDatabaseError, err)
	}
	d.db = db

	err = db.AutoMigrate(
		&postgresStorageORM.UserPG{},
		&postgresStorageORM.ToDoItemPG{},
	)

	if err != nil {
		return errors.Join(storage.ErrDatabaseError, err)
	}

	return nil
}

func (d *PostgresDataProvider) Run() error {
	return d.runWithDialector(postgres.Open(d.dsn), true)
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
func (d *PostgresDataProvider) StorageToDoItemCreate(ctx context.Context, title string, ownerID uint64) (uint64, error) {
	newItem := postgresStorageORM.ToDoItemPG{
		OwnerID: ownerID,
		Title:   title,
	}

	result := d.db.WithContext(ctx).Create(&newItem)
	if result.Error != nil {
		return 0, errors.Join(storage.ErrDatabaseError, result.Error)
	}

	return newItem.ID, nil
}

// StorageToDoItem_Update implements todoService.IToDoItemUpdater.
func (d *PostgresDataProvider) StorageToDoItemUpdate(ctx context.Context, item storageDTO.ToDoItem, ownerID uint64) error {

	updatedFields := make(map[string]interface{}, 2)

	if item.Title != nil {
		updatedFields["title"] = *item.Title
	}

	if item.IsComplete != nil {
		updatedFields["is_complete"] = *item.IsComplete
	}

	result := d.db.WithContext(ctx).Model(&postgresStorageORM.ToDoItemPG{ID: item.Id}).Updates(updatedFields)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return storage.ErrNotFound
		}
		errors.Join(storage.ErrDatabaseError, result.Error)
	}

	return nil
}

// StorageToDoItem_GetById implements todoService.IToDoItemGetter.
func (d *PostgresDataProvider) StorageToDoItemGetByID(ctx context.Context, itemID uint64, ownerID uint64) (*storageDTO.ToDoItem, error) {
	var item postgresStorageORM.ToDoItemPG
	result := d.db.WithContext(ctx).First(&item, itemID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, storage.ErrNotFound
		}
		return nil, errors.Join(storage.ErrDatabaseError, result.Error)
	}

	return &storageDTO.ToDoItem{
		Id:         item.ID,
		Title:      &item.Title,
		IsComplete: &item.IsComplete,
		OwnerId:    item.OwnerID,
	}, nil
}

// StorageToDoItem_GetList implements todoService.IToDoItemGetter.
func (d *PostgresDataProvider) StorageToDoItemGetList(ctx context.Context, ownerID uint64) ([]storageDTO.ToDoItem, error) {
	var items []postgresStorageORM.ToDoItemPG
	var resultList []storageDTO.ToDoItem

	result := d.db.WithContext(ctx).Find(&items, "owner_id = ?", ownerID)
	if result.Error != nil {
		return resultList, result.Error
	}

	for _, item := range items {
		resultList = append(resultList, storageDTO.ToDoItem{
			Id:         item.ID,
			Title:      &item.Title,
			IsComplete: &item.IsComplete,
			OwnerId:    item.OwnerID,
		})
	}

	return resultList, nil
}

// StorageToDoItem_DeleteById implements todoService.IToDoItemDeleter.
func (d *PostgresDataProvider) StorageToDoItemDeleteByID(ctx context.Context, itemID uint64, ownerID uint64) error {
	result := d.db.WithContext(ctx).Delete(&postgresStorageORM.ToDoItemPG{}, itemID)

	if result.Error != nil {
		return errors.Join(storage.ErrDatabaseError, result.Error)
	}

	if result.RowsAffected == 0 {
		return storage.ErrNotFound
	}

	return nil
}

// GetCredential implements credentialService.ICredentialStorageProvider.
func (d *PostgresDataProvider) GetCredential(username string) (*storageDTO.Credential, error) {
	log := d.log.With(slog.String("method", "GetCredential"))
	log.Warn("get credential not implement")
	return &storageDTO.Credential{
		Username:  "user",
		TokenHash: nil,
	}, nil
}

// var _ authService.IAccountCreator = (*PostgresDataProvider)(nil)
// var _ authService.IAccountGetter = (*PostgresDataProvider)(nil)

// GetAccountById implements authService.IAccountGetter.
func (d *PostgresDataProvider) GetAccountByID(ctx context.Context, userID uint64) (*storageDTO.User, error) {
	var user storageDTO.User
	result := d.db.WithContext(ctx).First(&user, userID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, storage.ErrNotFound
		}
		return nil, errors.Join(storage.ErrDatabaseError, result.Error)
	}

	return &user, nil
}

// GetAccountByUsername implements authService.IAccountGetter.
func (d *PostgresDataProvider) GetAccountByUsername(ctx context.Context, username string) (*storageDTO.User, error) {
	var user storageDTO.User
	result := d.db.WithContext(ctx).First(&user, "username=?", username)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, storage.ErrNotFound
		}
		return nil, errors.Join(storage.ErrDatabaseError, result.Error)
	}

	return &user, nil
}

// CreateAccount implements authService.IAccountCreator.
func (d *PostgresDataProvider) CreateAccount(ctx context.Context, username string, passwordHash []byte) (*serviceDTO.User, error) {
	newUser := postgresStorageORM.UserPG{
		Username:     username,
		PasswordHash: passwordHash,
	}

	result := d.db.WithContext(ctx).Create(&newUser)
	if result.Error != nil {
		return nil, errors.Join(storage.ErrDatabaseError, result.Error)
	}

	return &serviceDTO.User{
		UserID:   &newUser.ID,
		Username: &newUser.Username,
	}, nil
}
