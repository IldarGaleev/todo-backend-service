// Package authservice implements authentication service
package authservice

import (
	"context"
	"errors"
	"log/slog"

	secretsDTO "github.com/IldarGaleev/todo-backend-service/internal/lib/secretsjwt/secretsdto"
	serviceDTO "github.com/IldarGaleev/todo-backend-service/internal/services/servicedto"
	"github.com/IldarGaleev/todo-backend-service/internal/storage"
	storageDTO "github.com/IldarGaleev/todo-backend-service/internal/storage/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrArguments   = errors.New("argument error")
	ErrNotFound    = errors.New("account not found")
	ErrWrongSecret = errors.New("wrong secret")
	ErrInternal    = errors.New("internal error")
)

//go:generate mockery --name IAccountGetter
type IAccountGetter interface {
	GetAccountByUsername(ctx context.Context, username string) (*storageDTO.User, error)
	GetAccountByID(ctx context.Context, userID uint64) (*storageDTO.User, error)
}

//go:generate mockery --name ISecretProvider
type ISecretProvider interface {
	CreateSecret(ctx context.Context, user secretsDTO.User) ([]byte, error)
	ValidateSecret(ctx context.Context, secret []byte) (*secretsDTO.User, error)
	DeleteSecret(ctx context.Context, secret []byte) error
}

type AuthService struct {
	logger         *slog.Logger
	accountGetter  IAccountGetter
	secretProvider ISecretProvider
}

func New(
	log *slog.Logger,
	secretProvider ISecretProvider,
	accountGetter IAccountGetter,
) *AuthService {
	return &AuthService{
		logger:         log.With(slog.String("module", "authService")),
		secretProvider: secretProvider,
		accountGetter:  accountGetter,
	}
}

func (s *AuthService) CheckSecret(ctx context.Context, secret []byte) (*serviceDTO.User, error) {
	log := s.logger.With(slog.String("method", "CheckSecret"))
	user, err := s.secretProvider.ValidateSecret(ctx, secret)
	if err != nil {
		log.Debug("wrong secret", slog.Any("err", err))
		return nil, ErrWrongSecret
	}
	return &serviceDTO.User{
		UserID:   user.UserID,
		Username: user.Username,
	}, nil
}

func (s *AuthService) DeleteSecret(ctx context.Context, secret []byte) error {
	log := s.logger.With(slog.String("method", "CheckSecret"))
	err := s.secretProvider.DeleteSecret(ctx, secret)
	if err != nil {
		log.Debug("wrong secret", slog.Any("err", err))
		return ErrWrongSecret
	}
	return nil
}

func (s *AuthService) CreateUserSecret(ctx context.Context, user serviceDTO.User) (string, error) {
	log := s.logger.With(slog.String("method", "CreateUserSecret"))

	if (user.UserID == nil && user.Username == nil) || user.Password == "" {
		log.Error("wrong arguments")
		return "", ErrArguments
	}

	var userAccount *storageDTO.User
	var err error

	if user.Username != nil {
		userAccount, err = s.accountGetter.GetAccountByUsername(ctx, *user.Username)
	} else {
		userAccount, err = s.accountGetter.GetAccountByID(ctx, *user.UserID)
	}

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", ErrNotFound
		}
		log.Error("get account error", slog.Any("err", err))
		return "", errors.Join(ErrInternal, err)
	}

	err = bcrypt.CompareHashAndPassword(userAccount.PasswordHash, []byte(user.Password))

	if err != nil {
		log.Debug("pasword hash compare error", slog.Any("err", err))
		return "", ErrWrongSecret
	}

	secretBytes, err := s.secretProvider.CreateSecret(ctx, secretsDTO.User{
		UserID:   &userAccount.Id,
		Username: &userAccount.Username,
	})

	if err != nil {
		log.Debug("pasword hash compare error", slog.Any("err", err))
		return "", errors.Join(ErrInternal, err)
	}

	return string(secretBytes), nil
}
