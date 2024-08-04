/*
Service entrypoint
*/

package app

import (
	"log/slog"

	configApp "github.com/IldarGaleev/todo-backend-service/internal/app/config"
	grpcApp "github.com/IldarGaleev/todo-backend-service/internal/app/grpc"
	credentialService "github.com/IldarGaleev/todo-backend-service/internal/services/credential"
	todoService "github.com/IldarGaleev/todo-backend-service/internal/services/todo"
	"github.com/IldarGaleev/todo-backend-service/internal/storage/fakedb"
)

type IStorageProvider interface {
	MustRun()
	Stop() error
}

// Main application
type App struct {
	grpcServer       *grpcApp.App
	todoItemsStorage IStorageProvider
}

// Create main application instance
func New(
	log *slog.Logger,
	config *configApp.AppConfig,
) *App {

	storageProvider := fakedb.New(log)

	return &App{
		grpcServer: grpcApp.New(
			log,
			config.Port,
			todoService.New(log, storageProvider),
			credentialService.New(log, storageProvider),
		),
		todoItemsStorage: storageProvider,
	}
}

func (app *App) MustRun() {
	app.grpcServer.MustRun()
	app.todoItemsStorage.MustRun()
}

func (app *App) Stop() {
	app.grpcServer.Stop()
	app.todoItemsStorage.Stop()
}
