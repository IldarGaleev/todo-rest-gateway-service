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

type MainApp struct {
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
) *MainApp {

	gRPCApp := grpcapplication.New(
		logger,
	)

	return &MainApp{
		logger:      logger,
		confApp:     appConf,
		grpcApp:     gRPCApp,
		apiBasePath: apiBasePath,
	}
}

func (rApp *MainApp) MustRun() {

	client, err := rApp.grpcApp.Start(rApp.confApp.Grpc.Hostname, rApp.confApp.Grpc.Port)

	if err != nil {
		panic(err)
	}

	authProvider := authprovider.New(rApp.logger, *client)
	todoProvider := todoprovider.New(rApp.logger, *client)

	authHandle := authhandler.New(rApp.logger, authProvider)
	authMiddleware := jwtmiddleware.New(rApp.logger, authProvider)
	todoItemHandler := todoitemshandler.New(
		rApp.logger,
		todoProvider,
		todoProvider,
		todoProvider,
		todoProvider,
	)

	httpApp := httpapplication.New(
		rApp.logger,
		rApp.apiBasePath,
		todoItemHandler,
		todoItemHandler,
		todoItemHandler,
		todoItemHandler,
		authHandle,
		authMiddleware,
	)

	rApp.httpApp = httpApp

	err = httpApp.Run(rApp.confApp.Api.Hostname, rApp.confApp.Api.Port)
	if err != nil {
		panic(err)
	}
}

func (rApp *MainApp) MustStop(ctx context.Context) {
	errHttp := rApp.httpApp.Stop(ctx)
	errGrpc := rApp.grpcApp.Stop()

	if errHttp != nil || errGrpc != nil {
		panic(errors.Join(ErrAppFailedStopServices, errHttp, errGrpc))
	}
}
