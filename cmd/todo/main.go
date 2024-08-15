package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todoapiservice/internal/app"
	"todoapiservice/internal/app/configapplication"
)

func main() {

	appConf := configapplication.MustLoadConfig()

	log := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	const apiBasePath = "/api/v1/"

	mainApp := app.New(
		log,
		appConf,
		apiBasePath,
	)

	mainApp.MustRun()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mainApp.MustStop(ctx)
	select {
	case <-ctx.Done():
		log.Warn("timeout of 5 seconds.")
	}
	log.Info("Server exiting")
}
