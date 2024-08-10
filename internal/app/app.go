/*
Package app Service entrypoint
*/
package app

import (
	"log/slog"

	configApp "github.com/IldarGaleev/todo-backend-service/internal/app/configapp"
	grpcApp "github.com/IldarGaleev/todo-backend-service/internal/app/grpcapp"
	secretsJwt "github.com/IldarGaleev/todo-backend-service/internal/lib/secretsjwt"
	authService "github.com/IldarGaleev/todo-backend-service/internal/services/auth"
	credentialService "github.com/IldarGaleev/todo-backend-service/internal/services/credentialservice"
	todoService "github.com/IldarGaleev/todo-backend-service/internal/services/todoservice"
	"github.com/IldarGaleev/todo-backend-service/internal/storage/postgresdb"
	faketempdb "github.com/IldarGaleev/todo-backend-service/internal/tempstorage/fakeTempDb"
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
	tokenStorage := faketempdb.New(log)

	secretProvider := secretsJwt.New(
		log,
		*config, tokenStorage,
		tokenStorage,
	)

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
			authService,
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
