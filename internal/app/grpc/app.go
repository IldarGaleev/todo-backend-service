/*
Implements gRPC Application
*/
package grpcApp

import (
	"fmt"
	"log/slog"
	"net"

	grpcToDoServer "github.com/IldarGaleev/todo-backend-service/internal/grpc/todo"
	"google.golang.org/grpc"
)

// gRPC Application
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// Create gRPC application instance
func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()
	grpcToDoServer.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// Run gRPC server listener, panic if failed
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run gRPC server listener
func (a *App) Run() error {
	const op = "grpcApp.Run"
	logger := a.log.With(
		slog.String("op", op),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info(
		"gRPC server started",
		slog.String("addr", listener.Addr().String()),
		slog.Int("port", a.port),
	)

	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop gRPC server listener
func (a *App) Stop() {
	const op = "grpcApp.Stop"

	logger := a.log.With(
		slog.String("op", op),
	)

	logger.Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
