/*
Implements gRPC Application
*/
package grpcApp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	grpcToDoServer "github.com/IldarGaleev/todo-backend-service/internal/grpc/todo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ICredentialService interface {
	CheckToken(token string) bool
}

// gRPC Application
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func GetUnaryInterceptor(credentialService ICredentialService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		meta, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing context metadata")
		}

		if len(meta["authorization"]) != 1 {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		if !credentialService.CheckToken(meta["authorization"][0]) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}

// Create gRPC application instance
func New(log *slog.Logger, port int, todoItemService grpcToDoServer.IToDoItemService, credentialSevice ICredentialService) *App {

	var opts []grpc.ServerOption

	opts = append(opts, grpc.UnaryInterceptor(GetUnaryInterceptor(credentialSevice)))

	//TODO: add TLS transport
	log.Warn("insecure transport for gRPC")
	gRPCServer := grpc.NewServer(opts...)

	grpcToDoServer.Register(gRPCServer, todoItemService)

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
