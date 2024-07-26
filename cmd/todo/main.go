package main

import (
	"log/slog"
	"os"

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

	slog.Info(
		"service start",
		slog.String("configString", config.AppConfig.ConfigStr),
	)
}
