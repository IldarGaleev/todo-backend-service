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
		logger:  log.With(slog.String("module","credentialService")),
		storage: storage,
	}
}

func (s *CredentialService) CheckToken(token string) bool {
	log:=s.logger.With(slog.String("method","CheckToken"))

	_, err := s.storage.GetCredential("123")
	if err != nil {
		log.Error("find credential error")
		return false
	}

	//TODO: unimplement token checker
	log.Error("unimplement credential token checker")
	return token == "Bearer 1234"
}
