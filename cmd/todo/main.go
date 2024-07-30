package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IldarGaleev/todo-backend-service/internal/app"
	configApp "github.com/IldarGaleev/todo-backend-service/internal/app/config"
	loggingApp "github.com/IldarGaleev/todo-backend-service/internal/app/logging"
)

func main() {

	//Init app config
	appConf := configApp.MustLoadConfig()

	//Init app logging
	log := loggingApp.New(
		loggingApp.EnvMode(appConf.EnvMode),
	)
	slog.SetDefault(log.Logging)

	//Init gRPC server
	grpcApp := app.New(
		log.Logging,
		appConf.Port,
	)

	go grpcApp.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop

	grpcApp.GRPCServer.Stop()

	slog.Info("application stopped", slog.String("signal", sig.String()))

}
