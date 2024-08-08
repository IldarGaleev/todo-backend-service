/*
Service entrypoint
*/

package app

import (
	"log/slog"

	configApp "github.com/IldarGaleev/todo-backend-service/internal/app/config"
	grpcApp "github.com/IldarGaleev/todo-backend-service/internal/app/grpc"
	secretsJwt "github.com/IldarGaleev/todo-backend-service/internal/lib/secrets"
	authService "github.com/IldarGaleev/todo-backend-service/internal/services/auth"
	credentialService "github.com/IldarGaleev/todo-backend-service/internal/services/credential"
	todoService "github.com/IldarGaleev/todo-backend-service/internal/services/todo"
	"github.com/IldarGaleev/todo-backend-service/internal/storage/postgresdb"
)

type IStorageProvider interface {
	MustRun()
	Stop() error
}

// Main application
type App struct {
	grpcServer      *grpcApp.App
	storageProvider IStorageProvider
}

// Create main application instance
func New(
	log *slog.Logger,
	config *configApp.AppConfig,
) *App {

	storageProvider := postgresdb.New(log, config.Dsn)

	secretProvider := secretsJwt.New(log,*config)

	todoService := todoService.New(
		log,
		storageProvider,
		storageProvider,
		storageProvider,
		storageProvider,
	)

	authService := authService.New(
		log,
		secretProvider,
		storageProvider,
	)

	return &App{
		grpcServer: grpcApp.New(
			log,
			config.Port,
			todoService,
			todoService,
			todoService,
			todoService,
			authService,
			credentialService.New(log, storageProvider),
		),
		storageProvider: storageProvider,
	}
}

func (app *App) MustRun() {
	app.storageProvider.MustRun()
	app.grpcServer.MustRun()
}

func (app *App) Stop() {
	app.grpcServer.Stop()
	app.storageProvider.Stop()
}
