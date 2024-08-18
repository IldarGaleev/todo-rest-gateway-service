package authprovider

import (
	"context"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
	"todoapiservice/internal/app/grpcapplication/mocks"
)

func TestAuthProvider_CheckSecret_ValidToken(t *testing.T) {
	logger := slog.Default()
	pr := mocks.New(false)
	instance := New(logger, pr)
	ctx := context.TODO()

	user, err := instance.CheckSecret(ctx, "1:user1")

	require.NoError(t, err)
	require.Equal(t, uint64(1), *user.UserID)
	require.Equal(t, "user1", *user.EMail)
}

func TestAuthProvider_CheckSecret_InvalidToken(t *testing.T) {
	logger := slog.Default()
	pr := mocks.New(false)
	instance := New(logger, pr)
	ctx := context.TODO()

	testData := []struct {
		name  string
		token string
	}{
		{
			name:  "Empty token",
			token: "",
		},
		{
			name:  "Valid token wrong data",
			token: "1:wronguser",
		},
		{
			name:  "Invalid token data",
			token: "1:user1:data",
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			user, err := instance.CheckSecret(ctx, data.token)

			require.ErrorIs(t, err, ErrPermissionDenied)
			require.Nil(t, user)
		})
	}
}

func TestAuthProvider_CheckSecret_ServerUnavailable(t *testing.T) {
	logger := slog.Default()
	pr := mocks.New(true)
	instance := New(logger, pr)
	ctx := context.TODO()

	user, err := instance.CheckSecret(ctx, "1:user1")

	require.ErrorIs(t, err, ErrAuthInternal)
	require.Nil(t, user)
}

func TestAuthProvider_Login_Valid(t *testing.T) {
	logger := slog.Default()
	pr := mocks.New(false)
	instance := New(logger, pr)
	ctx := context.TODO()

	user, err := instance.Login(ctx, "user1", "pass")

	require.NoError(t, err)
	require.Equal(t, "1:user1", *user.JWT)
}

func TestAuthProvider_Login_Invalid(t *testing.T) {
	logger := slog.Default()
	pr := mocks.New(false)
	instance := New(logger, pr)
	ctx := context.TODO()

	testData := []struct {
		name     string
		eMail    string
		password string
	}{
		{
			name:     "Empty data",
			eMail:    "",
			password: "",
		},
		{
			name:     "Empty username",
			eMail:    "",
			password: "pass",
		},
		{
			name:     "Empty password",
			eMail:    "user1",
			password: "",
		},
		{
			name:     "Invalid password",
			eMail:    "user1",
			password: "invalid",
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			user, err := instance.Login(ctx, data.eMail, data.password)

			require.ErrorIs(t, err, ErrPermissionDenied)
			require.Nil(t, user)
		})
	}
}

func TestAuthProvider_Login_ServerUnavailable(t *testing.T) {
	logger := slog.Default()
	pr := mocks.New(true)
	instance := New(logger, pr)
	ctx := context.TODO()

	user, err := instance.Login(ctx, "user1", "pass")

	require.ErrorIs(t, err, ErrAuthInternal)
	require.Nil(t, user)
}
