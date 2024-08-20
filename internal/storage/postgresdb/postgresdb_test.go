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

	rows := sqlmock.NewRows(
		[]string{
			"id",
			"username",
			"password_hash",
		}).AddRow(
		1,
		"user",
		[]byte("hash"),
	)

	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	usr, err := storageService.GetAccountByID(ctx, 1)

	require.NoError(t, err)

	require.Equal(t, "user", usr.Username)
	require.Equal(t, []byte("hash"), usr.PasswordHash)
	require.Equal(t, uint64(1), usr.Id)
}

func TestPostgresDataProvider_GetAccountByID_Error(t *testing.T) {
	ctx := context.Background()
	storageService, mock := createStorage(t)
	mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrRecordNotFound)

	usr, err := storageService.GetAccountByID(ctx, 1)
	require.ErrorIs(t, err, storage.ErrNotFound)
	require.Nil(t, usr)
}
