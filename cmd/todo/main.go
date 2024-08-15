package main

import (
	"context"
	"log/slog"
	"os"
	"todoapiservice/internal/app/httpapplication"
	"todoapiservice/internal/http/handlers/authhandler"
	"todoapiservice/internal/http/handlers/todoitemshandler"
	"todoapiservice/internal/http/middlewares/jwtmiddleware"
	"todoapiservice/internal/services/authprovider"
	"todoapiservice/internal/services/todoprovider"

	todoprotobufv1 "github.com/IldarGaleev/todo-backend-service/pkg/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func unaryInterceptor(
	ctx context.Context,
	method string, req, reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {

	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+"1234")
	err := invoker(authCtx, method, req, reply, cc, opts...)

	// If we got an unauthenticated response from the gRPC service, refresh the token
	if status.Code(err) == codes.Unauthenticated {
		// if err = jwt.refreshBearerToken(); err != nil {
		//     return err
		// }
		updatedAuthCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+"1234")
		err = invoker(updatedAuthCtx, method, req, reply, cc, opts...)
	}

	return err
}

func main() {

	log := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithUnaryInterceptor(unaryInterceptor))

	conn, err := grpc.NewClient("localhost:9090", opts...)

	if err != nil {
		panic(err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	client := todoprotobufv1.NewToDoServiceClient(conn)

	authProvider := authprovider.New(log, client)
	todoProvider := todoprovider.New(log, client)

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

	err = httpApp.Run("localhost", 8080)

	if err != nil {
		panic(err)
	}
}
