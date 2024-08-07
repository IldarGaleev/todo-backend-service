package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IldarGaleev/todo-backend-service/internal/app"
	configApp "github.com/IldarGaleev/todo-backend-service/internal/app/config"
	appLogging "github.com/IldarGaleev/todo-backend-service/internal/lib/logging"
)

func main() {

	//Init app config
	appConf := configApp.MustLoadConfig()

	//Init app logging
	log := appLogging.New(
		appLogging.EnvMode(appConf.EnvMode),
	)
	slog.SetDefault(log.Logging)

	//Init gRPC server
	grpcApp := app.New(
		log.Logging,
		appConf,
	)

	go grpcApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop

	grpcApp.Stop()

	slog.Info("application stopped", slog.String("signal", sig.String()))

}
