// Package secretsjwt implements secrets for JWT
package secretsjwt

import (
	"context"
	"errors"
	"log/slog"
	"time"

	configApp "github.com/IldarGaleev/todo-backend-service/internal/app/configapp"
	secretsDTO "github.com/IldarGaleev/todo-backend-service/internal/lib/secretsjwt/secretsdto"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrVerifyError = errors.New("jwt verify error")
	ErrCreateError = errors.New("jwt create error")
)

type IJWTIndexer interface {
	CreateNewID(ctx context.Context) uint64
}

type IJWTRevoker interface {
	IsJWTRevoked(ctx context.Context, id uint64) bool
	RevokeJWT(ctx context.Context, id uint64)
}

type SecretJWT struct {
	logger     *slog.Logger
	maxAge     time.Duration
	secretKey  []byte
	jwtIndexer IJWTIndexer
	jwtRevoker IJWTRevoker
}

type TokenClaims struct {
	UserID   uint64 `json:"userid"`
	Username string `json:"username"`
	TokenID  uint64 `json:"tokenid"`
	jwt.RegisteredClaims
}

func New(
	log *slog.Logger,
	config configApp.AppConfig,
	jwtIndexer IJWTIndexer,
	jwtRevoker IJWTRevoker,
) *SecretJWT {
	return &SecretJWT{
		logger:     log.With(slog.String("module", "secretsJwt")),
		maxAge:     config.SecretsMaxAge,
		secretKey:  config.SecretKey,
		jwtIndexer: jwtIndexer,
		jwtRevoker: jwtRevoker,
	}
}

func (s *SecretJWT) decodeToken(secret []byte) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		string(secret),
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return s.secretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok {
		return claims, nil
	}

	return nil, ErrVerifyError
}

func (s *SecretJWT) ValidateSecret(ctx context.Context, secret []byte) (*secretsDTO.User, error) {
	log := s.logger.With(slog.String("method", "ValidateSecret"))

	claims, err := s.decodeToken(secret)
	if err != nil {
		log.Debug("jwt parse error", slog.Any("err", err))
		return nil, ErrVerifyError
	}

	if s.jwtRevoker.IsJWTRevoked(ctx, claims.TokenID) {
		return nil, ErrVerifyError
	}

	return &secretsDTO.User{
		UserID:   &claims.UserID,
		Username: &claims.Username,
	}, nil

}

func (s *SecretJWT) CreateSecret(ctx context.Context, user secretsDTO.User) ([]byte, error) {
	log := s.logger.With(slog.String("method", "CreateSecret"))

	claims := TokenClaims{
		UserID:   *user.UserID,
		Username: *user.Username,
		TokenID:  s.jwtIndexer.CreateNewID(ctx),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.maxAge)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "user",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secretKey)

	if err != nil {
		log.Debug("jwt create error", slog.Any("err", err))
		return nil, ErrCreateError
	}

	return []byte(tokenString), nil
}

func (s *SecretJWT) DeleteSecret(ctx context.Context, secret []byte) error {
	log := s.logger.With(slog.String("method", "DeleteSecret"))

	claims, err := s.decodeToken(secret)

	if err != nil {
		log.Debug("jwt parse error", slog.Any("err", err))
		return ErrVerifyError
	}
	if s.jwtRevoker.IsJWTRevoked(ctx, claims.TokenID) {
		return ErrVerifyError
	}
	s.jwtRevoker.RevokeJWT(ctx, claims.TokenID)
	return nil
}
