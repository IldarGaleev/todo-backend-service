package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IldarGaleev/todo-backend-service/internal/app"
	"github.com/IldarGaleev/todo-backend-service/internal/config"
)

func initLogging() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	return nil
}

func main() {

	err := initLogging()

	if err != nil {
		panic(err)
	}

	err = config.LoadConfig()
	if err != nil {
		slog.Error(
			"config load",
			slog.String("error", err.Error()),
		)
	}

	application := app.New(slog.Default(), config.AppConfig.Port)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop

	application.GRPCServer.Stop()

	slog.Info("application stopped", slog.String("signal", sig.String()))

}
