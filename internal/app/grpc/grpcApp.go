/*
Implements gRPC Application
*/
package grpcApp

import (
	"context"
	"errors"
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

var (
	ErrGrpcServe  = errors.New("grpc app: serve error")
	ErrGrpcListen = errors.New("grpc app: listen error")
)

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
func New(
	log *slog.Logger,
	port int,
	todoItemsCreatorService grpcToDoServer.IToDoItemCreatorService,
	todoItemsUpdaterService grpcToDoServer.IToDoItemUpdaterService,
	todoItemsGetterService grpcToDoServer.IToDoItemGetterService,
	todoItemsDeleterService grpcToDoServer.IToDoItemDeleterService,
	accountSecretCreator grpcToDoServer.IAccountSecretCreator,
	accountSecretValidator grpcToDoServer.IAccountSecretValidator,
	accountSecretDeleter grpcToDoServer.IAccountSecretDeleter,
	credentialSevice ICredentialService,
) *App {

	var opts []grpc.ServerOption

	opts = append(opts, grpc.UnaryInterceptor(GetUnaryInterceptor(credentialSevice)))

	//TODO: add TLS transport
	log.Warn("insecure transport for gRPC")
	gRPCServer := grpc.NewServer(opts...)

	grpcToDoServer.Register(
		gRPCServer,
		todoItemsCreatorService,
		todoItemsUpdaterService,
		todoItemsGetterService,
		todoItemsDeleterService,
		accountSecretCreator,
		accountSecretValidator,
		accountSecretDeleter,
	)

	return &App{
		log:        log.With(slog.String("module", "grpcApp")),
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
	log := a.log.With(slog.String("method", "Run"))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))

	if err != nil {
		return errors.Join(ErrGrpcListen, err)
	}

	log.Info(
		"gRPC server started",
		slog.String("addr", listener.Addr().String()),
		slog.Int("port", a.port),
	)

	if err := a.gRPCServer.Serve(listener); err != nil {
		return errors.Join(ErrGrpcServe, err)
	}

	return nil
}

// Stop gRPC server listener
func (a *App) Stop() {
	log := a.log.With(slog.String("method", "Stop"))

	log.Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
