package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todoapiservice/internal/app"
	"todoapiservice/internal/app/configapplication"
	"todoapiservice/internal/lib/applogging"
)

func main() {

	confPath := "config.yml"

	appConf := configapplication.MustLoadConfig(confPath)

	loggingApp := applogging.New(applogging.EnvMode(appConf.EnvMode))

	logging := loggingApp.Logging.With("module", "main")

	const apiBasePath = "/api/v1/"

	mainApp := app.New(
		loggingApp.Logging,
		appConf,
		apiBasePath,
	)

	go mainApp.MustRun()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mainApp.MustStop(ctx)

	<-ctx.Done()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		logging.Warn("app stop timeout")
	}
	logging.Info("Server exiting")
}
