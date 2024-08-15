package main

import (
	"log/slog"
	"os"
	"todoapiservice/internal/app/configapplication"
	"todoapiservice/internal/app/grpcapplication"
	"todoapiservice/internal/app/httpapplication"
	"todoapiservice/internal/http/handlers/authhandler"
	"todoapiservice/internal/http/handlers/todoitemshandler"
	"todoapiservice/internal/http/middlewares/jwtmiddleware"
	"todoapiservice/internal/services/authprovider"
	"todoapiservice/internal/services/todoprovider"
)

func main() {

	log := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	appConf := configapplication.MustLoadConfig()

	gRPCApp := grpcapplication.New(
		log,
	)

	client, err := gRPCApp.Start(appConf.GrpcHostname, appConf.GrpcPort)

	if err != nil {
		panic(err)
	}

	defer func() {
		err := gRPCApp.Stop()
		if err != nil {
			panic(err)
		}
	}()

	authProvider := authprovider.New(log, *client)
	todoProvider := todoprovider.New(log, *client)

	authHandle := authhandler.New(log, authProvider)
	authMiddleware := jwtmiddleware.New(log, authProvider)
	todoItemHandler := todoitemshandler.New(
		log,
		todoProvider,
		todoProvider,
		todoProvider,
		todoProvider,
	)

	const apiBasePath = "/api/v1/"

	httpApp := httpapplication.New(
		log,
		apiBasePath,
		todoItemHandler,
		todoItemHandler,
		todoItemHandler,
		todoItemHandler,
		authHandle,
		authMiddleware,
	)

	err = httpApp.Run(appConf.APIHostname, appConf.APIPort)

	if err != nil {
		panic(err)
	}
}
