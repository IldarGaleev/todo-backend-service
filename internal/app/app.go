/*
Service entrypoint
*/

package app

import (
	"log/slog"

	configApp "github.com/IldarGaleev/todo-backend-service/internal/app/config"
	grpcApp "github.com/IldarGaleev/todo-backend-service/internal/app/grpc"
	todoService "github.com/IldarGaleev/todo-backend-service/internal/services/todo"
	"github.com/IldarGaleev/todo-backend-service/internal/storage/fakedb"
)

// Main application
type App struct {
	grpcServer *grpcApp.App
	// todoItemsService *todoService.TodoService
	storageProvider *fakedb.FakeDatabaseProvider
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
		),
		storageProvider: storageProvider,
	}
}

func (app *App) MustRun() {
	app.grpcServer.MustRun()
	app.storageProvider.MustRun()
}

func (app *App) Stop() {
	app.grpcServer.Stop()
	app.storageProvider.Stop()
}
