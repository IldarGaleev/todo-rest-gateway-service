package grpcapplication

import (
	"context"
	"errors"
	"fmt"
	todoprotobufv1 "github.com/IldarGaleev/todo-backend-service/pkg/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
)

var (
	ErrGRPCStartError = errors.New("error starting gRPC server")
	ErrGRPCStopError  = errors.New("error stopping gRPC server")
	ErrGRPCNotRunning = errors.New("gRPC server is not running")
)

type GRPCApplication struct {
	logger *slog.Logger
	conn   *grpc.ClientConn
}

func unaryInterceptor(
	ctx context.Context,
	method string, req, reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	//TODO: unaryInterceptor not implement
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

func New(logger *slog.Logger) *GRPCApplication {
	log := logger.With("module", "grpcapplication")
	log.With("method", "New").Error("gRPC unaryInterceptor not implement!")
	return &GRPCApplication{
		logger: log,
	}
}

func (app *GRPCApplication) Start(host string, port int) (*todoprotobufv1.ToDoServiceClient, error) {
	log := app.logger.With("method", "Start")
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithUnaryInterceptor(unaryInterceptor))

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), opts...)

	if err != nil {
		log.Error("failed to create gRPC client", slog.Any("err", err))
		return nil, errors.Join(ErrGRPCStartError, err)
	}

	app.conn = conn

	client := todoprotobufv1.NewToDoServiceClient(conn)
	return &client, nil
}

func (app *GRPCApplication) Stop() error {
	log := app.logger.With("method", "Stop")
	if app.conn == nil {
		return ErrGRPCNotRunning
	}

	err := app.conn.Close()
	if err != nil {
		log.Error("failed to close gRPC connection", slog.Any("err", err))
		return errors.Join(ErrGRPCStopError, err)
	}

	return nil
}
