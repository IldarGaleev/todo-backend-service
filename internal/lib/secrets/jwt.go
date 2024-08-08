package secretsJwt

import (
	"context"
	"log/slog"
	"time"

	configApp "github.com/IldarGaleev/todo-backend-service/internal/app/config"
	secretsDTO "github.com/IldarGaleev/todo-backend-service/internal/lib/secrets/models"
	"github.com/kataras/jwt"
)

type SecretJWT struct {
	logger *slog.Logger
	maxAge time.Duration
	secretKey []byte
}

func New(log *slog.Logger, config configApp.AppConfig) (*SecretJWT){
	return &SecretJWT{
		logger: log,
		maxAge: config.SecretsMaxAge,
		secretKey: config.SecretKey,
	}
}

func (s *SecretJWT)ValidateSecret(ctx context.Context, secret []byte) (error){
	// log:=s.logger.With(slog.String("method","ValidateSecret"))
	return nil
}

func (s *SecretJWT)CreateSecret(ctx context.Context, user secretsDTO.User) ([]byte, error) {
	log:=s.logger.With(slog.String("method","CreateSecret"))
	token, err := jwt.Sign(jwt.HS256, s.secretKey, user.Payload, jwt.MaxAge(s.maxAge))
	if err != nil {
		log.Debug("jwt validation error",slog.Any("err",err))
		return nil,err
	}
	return token, nil
}