// Package app implements API service application init
package app

import (
	"context"
	"errors"
	todoprotobufv1 "github.com/IldarGaleev/todo-backend-service/pkg/grpc/proto"
	"log/slog"
	"todoapiservice/internal/app/configapplication"
	"todoapiservice/internal/app/grpcapplication"
	"todoapiservice/internal/app/httpapplication"
	"todoapiservice/internal/http/handlers/authhandler"
	"todoapiservice/internal/http/handlers/todoitemshandler"
	"todoapiservice/internal/http/middlewares/jwtmiddleware"
	"todoapiservice/internal/services/authprovider"
	"todoapiservice/internal/services/todoprovider"
)

var (
	ErrAppFailedStopServices = errors.New("app failed to stop services")
)

type IGRPCClient interface {
	Start(host string, port int) (*todoprotobufv1.ToDoServiceClient, error)
	Stop() error
}

type IHTTPServer interface {
	Run(host string, port int) error
	Stop(ctx context.Context) error
}

type App struct {
	logger      *slog.Logger
	confApp     *configapplication.AppConfig
	grpcApp     IGRPCClient
	httpApp     IHTTPServer
	apiBasePath string
}

func New(
	logger *slog.Logger,
	appConf *configapplication.AppConfig,
	apiBasePath string,
) *App {

	gRPCApp := grpcapplication.New(
		logger,
	)

	return &App{
		logger:      logger,
		confApp:     appConf,
		grpcApp:     gRPCApp,
		apiBasePath: apiBasePath,
	}
}

func (app *App) MustRun() {

	client, err := app.grpcApp.Start(app.confApp.GrpcHostname, app.confApp.GrpcPort)

	if err != nil {
		panic(err)
	}

	authProvider := authprovider.New(app.logger, *client)
	todoProvider := todoprovider.New(app.logger, *client)

	authHandle := authhandler.New(app.logger, authProvider)
	authMiddleware := jwtmiddleware.New(app.logger, authProvider)
	todoItemHandler := todoitemshandler.New(
		app.logger,
		todoProvider,
		todoProvider,
		todoProvider,
		todoProvider,
	)

	httpApp := httpapplication.New(
		app.logger,
		app.apiBasePath,
		todoItemHandler,
		todoItemHandler,
		todoItemHandler,
		todoItemHandler,
		authHandle,
		authMiddleware,
	)

	err = httpApp.Run(app.confApp.APIHostname, app.confApp.APIPort)
	if err != nil {
		panic(err)
	}

	app.httpApp = httpApp
}

func (app *App) MustStop(ctx context.Context) {
	errHttp := app.httpApp.Stop(ctx)
	errGrpc := app.grpcApp.Stop()

	if errHttp != nil || errGrpc != nil {
		panic(errors.Join(ErrAppFailedStopServices, errHttp, errGrpc))
	}
}
