// Package authhandler implements authentication http handlers
package authhandler

import (
	"context"
	"log/slog"
	"net/http"

	"todoapiservice/internal/http/handlers"
	"todoapiservice/internal/http/httpdto"
	"todoapiservice/internal/services/coredto"

	"github.com/gin-gonic/gin"
)

type IAuthenticator interface {
	Login(ctx context.Context, email string, password string) (*coredto.User, error)
	Logout(ctx context.Context, user coredto.User) error
}

type AuthHandler struct {
	logging       *slog.Logger
	authenticator IAuthenticator
}

func New(
	logging *slog.Logger,
	authenticator IAuthenticator,
) *AuthHandler {
	return &AuthHandler{
		logging:       logging.With("module", "authhandler"),
		authenticator: authenticator,
	}
}

func (h *AuthHandler) HandlerLogin(c *gin.Context) {

	email, pass, hasAuth := c.Request.BasicAuth()

	if !hasAuth {
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		handlers.SendErrorResponse(c, http.StatusUnauthorized)
		return
	}

	user, err := h.authenticator.Login(c.Request.Context(), email, pass)

	if err != nil {
		handlers.SendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, httpdto.LoginResponse{
		GeneralResponse: httpdto.GeneralResponse{
			Status: httpdto.StatusOK,
		},
		Token: *user.JWT,
	})
}

func (h *AuthHandler) HandlerLogout(c *gin.Context) {
	token := c.GetString("jwtToken")

	err := h.authenticator.Logout(
		c.Request.Context(),
		coredto.User{
			JWT: &token,
		})

	if err != nil {
		handlers.SendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, httpdto.GeneralResponse{
		Status: httpdto.StatusOK,
	})
}
