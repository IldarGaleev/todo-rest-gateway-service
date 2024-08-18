// Package authprovider implements authentication gRPC bindings
package authprovider

import (
	"context"
	"errors"
	"log/slog"
	"todoapiservice/internal/services/coredto"

	todoprotobufv1 "github.com/IldarGaleev/todo-backend-service/pkg/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrAuthInternal     = errors.New("authentication internal error")
	ErrPermissionDenied = errors.New("authentication internal error")
)

type AuthProvider struct {
	logger *slog.Logger
	client todoprotobufv1.ToDoServiceClient
}

func New(
	logger *slog.Logger,
	client todoprotobufv1.ToDoServiceClient,
) *AuthProvider {
	return &AuthProvider{
		logger: logger.With("module", "authprovider"),
		client: client,
	}
}

func (p *AuthProvider) Login(ctx context.Context, email string, password string) (*coredto.User, error) {
	log := p.logger.With("method", "Login")
	resp, err := p.client.Login(ctx, &todoprotobufv1.LoginRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return nil, ErrPermissionDenied
		}

		log.Error("login error", slog.Any("err", err))
		return nil, errors.Join(ErrAuthInternal, err)
	}

	return &coredto.User{
		JWT: &resp.Token,
	}, nil

}

func (p *AuthProvider) Logout(ctx context.Context, user coredto.User) error {
	log := p.logger.With("method", "Logout")
	_, err := p.client.Logout(ctx, &todoprotobufv1.LogoutRequest{
		Token: *user.JWT,
	})

	if err != nil {
		log.Error("logout error", slog.Any("err", err))
		return ErrAuthInternal
	}

	return nil

}

func (p *AuthProvider) CheckSecret(ctx context.Context, secret string) (*coredto.User, error) {
	// log := p.logger.With("method","CheckSecret")
	resp, err := p.client.CheckSecret(ctx, &todoprotobufv1.CheckSecretRequest{
		Secret: secret,
	})

	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			return nil, ErrPermissionDenied
		}
		if status.Code(err) == codes.InvalidArgument {
			return nil, ErrPermissionDenied
		}
		return nil, errors.Join(ErrAuthInternal, err)
	}

	userID := resp.GetUserId()
	email := resp.GetEmail()

	return &coredto.User{
		UserID: &userID,
		EMail:  &email,
		JWT:    &secret,
	}, nil
}
