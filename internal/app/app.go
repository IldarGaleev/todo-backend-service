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

// App Main application
type App struct {
	logger          *slog.Logger
	grpcServer      *grpcApp.App
	storageProvider IStorageProvider
}

// New Create main application instance
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

	todoSrv := todoService.New(
		log,
		storageProvider,
		storageProvider,
		storageProvider,
		storageProvider,
	)

	authSrv := authService.New(
		log,
		secretProvider,
		storageProvider,
	)

	return &App{
		logger: log.With("module", "app"),
		grpcServer: grpcApp.New(
			log,
			config.Port,
			todoSrv,
			todoSrv,
			todoSrv,
			todoSrv,
			authSrv,
			authSrv,
			authSrv,
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
	err := app.storageProvider.Stop()
	if err != nil {
		app.logger.Error("failed stop service", slog.Any("err", err))
	}
}
