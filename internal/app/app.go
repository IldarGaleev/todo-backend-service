/*
Service entrypoint
*/

package app

import (
	"log/slog"

	grpcApp "github.com/IldarGaleev/todo-backend-service/internal/app/grpc"
)

// Main application
type App struct {
	GRPCServer *grpcApp.App
}

// Create main application instance
func New(
	log *slog.Logger,
	grpcPort int,
) *App {

	grpcApp := grpcApp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
