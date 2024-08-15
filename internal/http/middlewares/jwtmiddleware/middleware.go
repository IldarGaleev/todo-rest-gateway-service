// Package jwtmiddleware implements JWT auth middleware
package jwtmiddleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"todoapiservice/internal/http/httpdto"
	"todoapiservice/internal/services/coredto"

	"github.com/gin-gonic/gin"
)

type ISecretChecker interface {
	CheckSecret(ctx context.Context, secret string) (*coredto.User, error)
}

type JWTMiddleware struct {
	loggger       *slog.Logger
	secretChecker ISecretChecker
}

func New(
	logger *slog.Logger,
	secretChecker ISecretChecker,
) *JWTMiddleware {
	return &JWTMiddleware{
		loggger:       logger.With("module", "jwtmiddleware"),
		secretChecker: secretChecker,
	}
}

func sendErrorStatus(c *gin.Context, code int) {
	c.Writer.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
	c.IndentedJSON(
		code,
		httpdto.GeneralResponse{
			Status: httpdto.StatusError,
		})
	c.Abort()
}

func (m *JWTMiddleware) CreateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			sendErrorStatus(c, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(token, " ")
		method, token := parts[0], parts[1]

		if strings.ToLower(method) != "bearer" {
			sendErrorStatus(c, http.StatusBadRequest)
			return
		}

		user, err := m.secretChecker.CheckSecret(
			c.Request.Context(),
			token,
		)

		if err != nil || user.JWT == nil || user.UserID == nil {
			sendErrorStatus(c, http.StatusUnauthorized)
			return
		}

		c.Set("userID", *user.UserID)
		c.Set("jwtToken", *user.JWT)
		c.Next()
	}
}
