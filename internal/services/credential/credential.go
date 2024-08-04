package credentialService

import (
	"log/slog"

	storageDTO "github.com/IldarGaleev/todo-backend-service/internal/storage/models"
)

type ICredentialStorageProvider interface {
	GetCredential(username string) (*storageDTO.Credential, error)
}

type CredentialService struct {
	logger  *slog.Logger
	storage ICredentialStorageProvider
}

func New(log *slog.Logger, storage ICredentialStorageProvider) *CredentialService {
	return &CredentialService{
		logger:  log,
		storage: storage,
	}
}

func (s *CredentialService) CheckToken(token string) bool {
	const op = "credentialService.CheckToken"
	logger := s.logger.With(slog.String("op", op))

	_, err := s.storage.GetCredential("123")
	if err != nil {
		logger.Error("find credential error")
		return false
	}

	//TODO: unimplement token checker
	logger.Error("unimplement credential token checker")
	return token == "Bearer 1234"
}
