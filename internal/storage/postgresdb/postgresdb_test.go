package postgresdb

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IldarGaleev/todo-backend-service/internal/storage"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"testing"
)

func createStorage(t *testing.T) (*PostgresDataProvider, sqlmock.Sqlmock) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockDb, mock, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	storageService := New(
		logger,
		"",
	)

	err := storageService.runWithDialector(dialector, true)

	//TODO: "all expectations were already fulfilled ...." https://github.com/go-gorm/gorm/issues/3565
	//require.NoError(t, err)
	_ = err

	//mock.ExpectQuery(
	//	`SELECT count(\*) FROM information_schema.tables`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	return storageService, mock
}

func TestPostgresDataProvider_GetAccountByID_Success(t *testing.T) {
	ctx := context.Background()

	storageService, mock := createStorage(t)

	userId := uint64(99)
	username := "user1"
	userPasswordHash := []byte("hash")

	rows := sqlmock.NewRows(
		[]string{
			"id",
			"username",
			"password_hash",
		}).AddRow(
		userId,
		username,
		userPasswordHash,
	)

	mock.ExpectQuery(`SELECT`).WithArgs(userId, 1).WillReturnRows(rows)

	usr, err := storageService.GetAccountByID(ctx, userId)

	require.NoError(t, err)

	require.Equal(t, username, usr.Username)
	require.Equal(t, userPasswordHash, usr.PasswordHash)
	require.Equal(t, userId, usr.Id)
}

func TestPostgresDataProvider_GetAccountByID_Error_NotFound(t *testing.T) {
	ctx := context.Background()
	storageService, mock := createStorage(t)

	userId := uint64(99)

	mock.ExpectQuery(`SELECT`).WithArgs(userId, 1).WillReturnError(gorm.ErrRecordNotFound)

	usr, err := storageService.GetAccountByID(ctx, userId)
	require.ErrorIs(t, err, storage.ErrNotFound)
	require.Nil(t, usr)
}

func TestPostgresDataProvider_GetAccountByID_Error_DBInternal(t *testing.T) {
	ctx := context.Background()
	storageService, mock := createStorage(t)

	userId := uint64(99)

	mock.ExpectQuery(`SELECT`).WithArgs(userId, 1).WillReturnError(gorm.ErrInvalidDB)

	usr, err := storageService.GetAccountByID(ctx, userId)
	require.ErrorIs(t, err, storage.ErrDatabaseError)
	require.Nil(t, usr)
}

func TestPostgresDataProvider_GetAccountByUsername_Success(t *testing.T) {
	ctx := context.Background()

	storageService, mock := createStorage(t)

	username := "user1"

	rows := sqlmock.NewRows(
		[]string{
			"id",
			"username",
			"password_hash",
		}).AddRow(
		1,
		username,
		[]byte("hash"),
	)

	mock.ExpectQuery(`SELECT`).WithArgs(username, 1).WillReturnRows(rows)

	usr, err := storageService.GetAccountByUsername(ctx, username)

	require.NoError(t, err)

	require.Equal(t, username, usr.Username)
	require.Equal(t, []byte("hash"), usr.PasswordHash)
	require.Equal(t, uint64(1), usr.Id)
}

func TestPostgresDataProvider_GetAccountByUsername_Error_NotFound(t *testing.T) {
	ctx := context.Background()
	storageService, mock := createStorage(t)

	username := "user1"

	mock.ExpectQuery(`SELECT`).WithArgs(username, 1).WillReturnError(gorm.ErrRecordNotFound)

	usr, err := storageService.GetAccountByUsername(ctx, username)
	require.ErrorIs(t, err, storage.ErrNotFound)
	require.Nil(t, usr)
}

func TestPostgresDataProvider_GetAccountByUsername_Error_DBInternal(t *testing.T) {
	ctx := context.Background()
	storageService, mock := createStorage(t)

	username := "user1"

	mock.ExpectQuery(`SELECT`).WithArgs(username, 1).WillReturnError(gorm.ErrInvalidDB)

	usr, err := storageService.GetAccountByUsername(ctx, username)
	require.ErrorIs(t, err, storage.ErrDatabaseError)
	require.Nil(t, usr)
}

func TestPostgresDataProvider_CreateAccount_Success(t *testing.T) {
	ctx := context.Background()
	storageService, mock := createStorage(t)

	username := "user1"
	passwordHash := []byte("pass")

	mock.ExpectBegin()
	mock.ExpectQuery(`^INSERT INTO "users" (.+)$`).
		WithArgs(username, passwordHash).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(1),
		)
	mock.ExpectCommit()

	usr, err := storageService.CreateAccount(ctx, username, passwordHash)
	require.NoError(t, err)
	require.Equal(t, "user1", *usr.Username)
}

func TestPostgresDataProvider_CreateAccount_Error_DBInternal(t *testing.T) {
	ctx := context.Background()
	storageService, mock := createStorage(t)

	username := "user1"
	passwordHash := []byte("pass")

	mock.ExpectBegin()
	mock.ExpectQuery(`^INSERT INTO "users" (.+)$`).
		WithArgs(
			username,
			passwordHash,
		).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	usr, err := storageService.CreateAccount(ctx, username, passwordHash)

	require.NoError(t, mock.ExpectationsWereMet())
	require.ErrorIs(t, err, storage.ErrDatabaseError)
	require.ErrorIs(t, err, gorm.ErrInvalidDB)
	require.Nil(t, usr)
}
