package httpapplication

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

var (
	ErrHttpAppRunError  = errors.New("http app run error")
	ErrHttpAppNotRun    = errors.New("http app not run")
	ErrHttpAppStopError = errors.New("http app stop error")
)

type IMiddleware interface {
	CreateMiddleware() gin.HandlerFunc
}

type IItemCreateHandler interface {
	CreateHandlerCreateTask() gin.HandlerFunc
}

type IItemGetterHandler interface {
	CreateHandlerGetTaskList() gin.HandlerFunc
	CreateHandlerGetTaskByID() gin.HandlerFunc
}

type IItemUpdateHandler interface {
	CreateHandlerUpdateTaskByID() gin.HandlerFunc
}

type IItemDeleteHandler interface {
	CreateHandlerDeleteTaskByID() gin.HandlerFunc
}

type IAuthHandler interface {
	CreateHandlerLogin() gin.HandlerFunc
	CreateHandlerLogout() gin.HandlerFunc
}

type HttpApp struct {
	logger *slog.Logger
	router *gin.Engine
	srv    *http.Server
}

func New(
	logger *slog.Logger,
	apiBasePath string,

	itemCreateHandler IItemCreateHandler,
	itemGetterHandler IItemGetterHandler,
	itemUpdateHandler IItemUpdateHandler,
	itemDeleteHandler IItemDeleteHandler,
	authHandler IAuthHandler,

	authMiddleware IMiddleware,

) *HttpApp {

	router := gin.Default()

	apiAuth := router.Group(apiBasePath)
	apiNoAuth := router.Group(apiBasePath)

	apiAuth.Use(authMiddleware.CreateMiddleware())

	apiAuth.POST("/tasks", itemCreateHandler.CreateHandlerCreateTask())
	apiAuth.GET("/tasks", itemGetterHandler.CreateHandlerGetTaskList())
	apiAuth.GET("/tasks/:id", itemGetterHandler.CreateHandlerGetTaskByID())
	apiAuth.PATCH("/tasks/:id", itemUpdateHandler.CreateHandlerUpdateTaskByID())
	apiAuth.DELETE("/tasks/:id", itemDeleteHandler.CreateHandlerDeleteTaskByID())
	apiAuth.GET("/logout", authHandler.CreateHandlerLogout())

	apiNoAuth.POST("/login", authHandler.CreateHandlerLogin())

	return &HttpApp{
		logger: logger.With(slog.String("module", "httpapplication")),
		router: router,
	}
}

func (app *HttpApp) Run(host string, port int) error {

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: app.router.Handler(),
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error("http app run error", slog.Any("err", err))
		return errors.Join(ErrHttpAppRunError, err)
	}

	return nil
}

func (app *HttpApp) Stop(ctx context.Context) error {
	if app.srv == nil {
		return ErrHttpAppNotRun
	}

	err := app.srv.Shutdown(ctx)
	if err != nil {
		app.logger.Error("http app stop error", slog.Any("err", err))
		return errors.Join(ErrHttpAppStopError, err)
	}

	return nil
}
