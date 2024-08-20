package authservice

import (
	"context"
	"errors"
	"github.com/IldarGaleev/todo-backend-service/internal/lib/secretsjwt/secretsdto"
	"github.com/IldarGaleev/todo-backend-service/internal/services/auth/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func createAuthService(t *testing.T) (*mocks.ISecretProvider, *mocks.IAccountGetter, *AuthService) {
	logger := slog.Default()

	secretProvider := mocks.NewISecretProvider(t)
	accountGetter := mocks.NewIAccountGetter(t)

	authService := New(
		logger,
		secretProvider,
		accountGetter,
	)

	return secretProvider, accountGetter, authService
}

func TestAuthService_CheckSecret_Valid(t *testing.T) {
	ctx := context.Background()
	secretProvider, _, authService := createAuthService(t)

	secret := []byte("secret")

	userId := uint64(1)
	username := "test_user"
	userSecret := secretsdto.User{
		UserID:   &userId,
		Username: &username,
	}

	secretProvider.On(
		"ValidateSecret",
		mock.Anything,
		secret,
	).Return(&userSecret, nil)

	usr, err := authService.CheckSecret(
		ctx,
		secret,
	)

	require.NoError(t, err)
	require.Equal(t, &userId, usr.UserID)
}

func TestAuthService_CheckSecret_Invalid(t *testing.T) {
	ctx := context.Background()
	secretProvider, _, authService := createAuthService(t)

	secretProvider.On(
		"ValidateSecret",
		mock.Anything,
		mock.Anything,
	).Return(nil, errors.New("invalid secret"))

	usr, err := authService.CheckSecret(
		ctx,
		[]byte("any_secret"),
	)

	require.ErrorIs(t, err, ErrWrongSecret)
	require.Nil(t, usr)
}
