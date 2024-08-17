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
	Middleware(c *gin.Context)
}

type IItemCreateHandler interface {
	HandlerCreateTask(c *gin.Context)
}

type IItemGetterHandler interface {
	HandlerGetTaskList(c *gin.Context)
	HandlerGetTaskByID(c *gin.Context)
}

type IItemUpdateHandler interface {
	HandlerUpdateTaskByID(c *gin.Context)
}

type IItemDeleteHandler interface {
	HandlerDeleteTaskByID(c *gin.Context)
}

type IAuthHandler interface {
	HandlerLogin(c *gin.Context)
	HandlerLogout(c *gin.Context)
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

	apiAuth.Use(authMiddleware.Middleware)

	apiAuth.POST("/tasks", itemCreateHandler.HandlerCreateTask)
	apiAuth.GET("/tasks", itemGetterHandler.HandlerGetTaskList)
	apiAuth.GET("/tasks/:id", itemGetterHandler.HandlerGetTaskByID)
	apiAuth.PATCH("/tasks/:id", itemUpdateHandler.HandlerUpdateTaskByID)
	apiAuth.DELETE("/tasks/:id", itemDeleteHandler.HandlerDeleteTaskByID)
	apiAuth.GET("/logout", authHandler.HandlerLogout)

	apiNoAuth.POST("/login", authHandler.HandlerLogin)

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
