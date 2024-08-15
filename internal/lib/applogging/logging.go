// Package applogging implements logging provider
package applogging

import (
	"log/slog"
	"os"
)

type EnvMode string

const (
	EnvModeLocal EnvMode = "local"
	EnvModeDev   EnvMode = "dev"
	EnvModeProd  EnvMode = "prod"
)

type LogApp struct {
	Logging *slog.Logger
}

// New Returns application logger
func New(mode EnvMode) *LogApp {
	var log *slog.Logger
	switch mode {
	case EnvModeLocal:
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
	case EnvModeDev:
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
