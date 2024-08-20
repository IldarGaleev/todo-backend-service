package authservice

import (
	"context"
	"errors"
	"github.com/IldarGaleev/todo-backend-service/internal/lib/secretsjwt/secretsdto"
	"github.com/IldarGaleev/todo-backend-service/internal/services/auth/mocks"
	"github.com/IldarGaleev/todo-backend-service/internal/services/servicedto"
	"github.com/IldarGaleev/todo-backend-service/internal/storage"
	storageDTO "github.com/IldarGaleev/todo-backend-service/internal/storage/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log/slog"
	"testing"
)

func createAuthService(t *testing.T) (*mocks.ISecretProvider, *mocks.IAccountGetter, *AuthService) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

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

func prepareAccountGetter(ag *mocks.IAccountGetter, userPassword string, t *testing.T) {

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)

	if err != nil {
		t.Fail()
		t.Logf("create hash error\n")
	}

	someUser := storageDTO.User{
		PasswordHash: pwdHash,
	}

	ag.On(
		"GetAccountByUsername",
		mock.Anything,
		mock.Anything,
	).Return(&someUser, nil)
}

func TestAuthService_CreateUserSecret_Valid(t *testing.T) {
	ctx := context.Background()
	secretProvider, accountGetter, authService := createAuthService(t)

	userPassword := "secret"
	username := "test_user"
	prepareAccountGetter(accountGetter, userPassword, t)

	secretProvider.On(
		"CreateSecret",
		mock.Anything,
		mock.Anything,
	).Return([]byte("generated_token"), nil)

	token, err := authService.CreateUserSecret(ctx, servicedto.User{Username: &username, Password: userPassword})

	require.NoError(t, err)
	require.Equal(t, "generated_token", token)
}

func TestAuthService_CreateUserSecret_CreateSecret_InternalError(t *testing.T) {
	ctx := context.Background()
	secretProvider, accountGetter, authService := createAuthService(t)

	userPassword := "secret"
	username := "test_user"
	prepareAccountGetter(accountGetter, userPassword, t)

	secretProvider.On(
		"CreateSecret",
		mock.Anything,
		mock.Anything,
	).Return(nil, errors.New("create secret internal error"))

	_, err := authService.CreateUserSecret(ctx, servicedto.User{Username: &username, Password: userPassword})

	require.ErrorIs(t, err, ErrInternal)
}

func TestAuthService_CreateUserSecret_ByUserID_NotFound(t *testing.T) {
	ctx := context.Background()
	_, accountGetter, authService := createAuthService(t)

	userId := uint64(1)

	accountGetter.On(
		"GetAccountByID",
		mock.Anything,
		mock.Anything,
	).Return(nil, storage.ErrNotFound)

	_, err := authService.CreateUserSecret(ctx, servicedto.User{UserID: &userId, Password: "pwd"})

	require.ErrorIs(t, err, ErrNotFound)
}

func TestAuthService_CreateUserSecret_ByUserID_AccountGetterInternalError(t *testing.T) {
	ctx := context.Background()
	_, accountGetter, authService := createAuthService(t)

	userId := uint64(1)

	accountGetter.On(
		"GetAccountByID",
		mock.Anything,
		mock.Anything,
	).Return(nil, errors.New("account getter internal error"))

	_, err := authService.CreateUserSecret(ctx, servicedto.User{UserID: &userId, Password: "pwd"})

	require.ErrorIs(t, err, ErrInternal)
}

func TestAuthService_CreateUserSecret_InvalidUserData(t *testing.T) {
	ctx := context.Background()
	_, _, authService := createAuthService(t)

	_, err := authService.CreateUserSecret(ctx, servicedto.User{})

	require.ErrorIs(t, err, ErrArguments)
}

func TestAuthService_CreateUserSecret_InvalidUserPassword(t *testing.T) {
	ctx := context.Background()
	_, accountGetter, authService := createAuthService(t)

	prepareAccountGetter(accountGetter, "secret", t)

	testCases := []struct {
		name          string
		password      string
		expectedError error
	}{
		{
			name:          "empty password",
			password:      "",
			expectedError: ErrArguments,
		},
		{
			name:          "invalid password",
			password:      "wrong_password",
			expectedError: ErrWrongSecret,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			username := "user"
			_, err := authService.CreateUserSecret(
				ctx,
				servicedto.User{
					Username: &username,
					Password: testCase.password,
				})

			require.ErrorIs(t, err, testCase.expectedError)
		})
	}
}

func TestAuthService_DeleteSecret_Success(t *testing.T) {
	ctx := context.Background()
	secretProvider, _, authService := createAuthService(t)

	secretProvider.On(
		"DeleteSecret",
		mock.Anything,
		mock.Anything,
	).Return(nil)

	err := authService.DeleteSecret(ctx, nil)

	require.NoError(t, err)
}

func TestAuthService_DeleteSecret_Failed(t *testing.T) {
	ctx := context.Background()
	secretProvider, _, authService := createAuthService(t)

	secretProvider.On(
		"DeleteSecret",
		mock.Anything,
		mock.Anything,
	).Return(errors.New("delete secret failed"))

	err := authService.DeleteSecret(ctx, nil)

	require.ErrorIs(t, err, ErrWrongSecret)
}
