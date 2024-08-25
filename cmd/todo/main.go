package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todoapiservice/docs"
	"todoapiservice/internal/app"
	"todoapiservice/internal/app/configapplication"
	"todoapiservice/internal/lib/applogging"
)

// @Title 			ToDo list app
// @Version 		1.0
// @Description 	Todo list API service
// @BasePath 		/api/v1/
// @Host			localhost:8080

// @License.name 	MIT
// @License.url 	https://mit-license.org/

// @Securitydefinitions.apikey 	ApiKeyAuth
// @In 							header
// @Name 						Authorization

// @securityDefinitions.basic 	BasicAuth
// @In 							header
// @Name 						Authorization
func main() {

	confPath := "config.yml"

	appConf := configapplication.MustLoadConfig(confPath)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", appConf.Api.Hostname, appConf.Api.Port)

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
