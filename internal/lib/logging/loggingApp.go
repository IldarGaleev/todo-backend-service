package appLogging

import (
	"log/slog"
	"os"
)

type EnvMode string

const (
	ENV_MODE_LOCAL EnvMode = "local"
	ENV_MODE_DEV   EnvMode = "dev"
	ENV_MODE_PROD  EnvMode = "prod"
)

type LogApp struct {
	Logging *slog.Logger
}

// Returns application logger
func New(mode EnvMode) *LogApp {
	var log *slog.Logger
	switch mode {
	case ENV_MODE_LOCAL:
		{
			log = slog.New(
				slog.NewTextHandler(
					os.Stdout,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			)
		}
	case ENV_MODE_DEV:
		{
			log = slog.New(
				slog.NewJSONHandler(
					os.Stdout,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			)
		}
	default:
		{
			log = slog.New(
				slog.NewTextHandler(
					os.Stdout,
					&slog.HandlerOptions{
						Level: slog.LevelWarn,
					},
				),
			)
		}
	}
	return &LogApp{
		Logging: log,
	}
}
