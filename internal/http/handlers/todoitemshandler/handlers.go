// Package todoitemshandler implements ToDo items http handlers
package todoitemshandler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"todoapiservice/internal/http/handlers"
	"todoapiservice/internal/http/httpdto"
	"todoapiservice/internal/services/coredto"
	"todoapiservice/internal/services/todoprovider"
)

type IToDoCreator interface {
	Create(ctx context.Context, owner coredto.User, title string) (*coredto.ToDoItem, error)
}

type IToDoDeleter interface {
	Delete(ctx context.Context, item coredto.ToDoItem) error
}

type IToDoGetter interface {
	GetByID(ctx context.Context, owner coredto.User, itemID uint64) (*coredto.ToDoItem, error)
	GetList(ctx context.Context, owner coredto.User) ([]coredto.ToDoItem, error)
}

type IToDoUpdater interface {
	Update(ctx context.Context, item coredto.ToDoItem) (*coredto.ToDoItem, error)
}

type ToDoHandlers struct {
	logging     *slog.Logger
	itemCreator IToDoCreator
	itemGetter  IToDoGetter
	itemUpdater IToDoUpdater
	itemDeleter IToDoDeleter
}

func New(
	logging *slog.Logger,
	itemCreator IToDoCreator,
	itemGetter IToDoGetter,
	itemUpdater IToDoUpdater,
	itemDeleter IToDoDeleter,
) *ToDoHandlers {
	return &ToDoHandlers{
		logging:     logging.With("module", "todoitemshandler"),
		itemCreator: itemCreator,
		itemGetter:  itemGetter,
		itemUpdater: itemUpdater,
		itemDeleter: itemDeleter,
	}
}

func (h *ToDoHandlers) HandlerCreateTask(c *gin.Context) {

	var changes httpdto.TaskItemChanges
	err := c.BindJSON(&changes)
	if err != nil {
		handlers.SendErrorResponse(c, http.StatusBadRequest)
		return
	}

	userID := c.GetUint64("userID")

	newItem, err := h.itemCreator.Create(
		c.Request.Context(),
		coredto.User{
			UserID: &userID,
		},
		*changes.Title,
	)

	if err != nil {
		handlers.SendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, httpdto.GetTaskByIDResponse{
		GeneralResponse: httpdto.GeneralResponse{
			Status: httpdto.StatusOK,
		},
		Task: httpdto.TaskItem{
			ID:    *newItem.ItemID,
			Title: *newItem.Title,
		},
	})

}

func (h *ToDoHandlers) HandlerGetTaskList(c *gin.Context) {
	userID := c.GetUint64("userID")

	items, err := h.itemGetter.GetList(
		c.Request.Context(),
		coredto.User{
			UserID: &userID,
		},
	)

	if err != nil {
		c.IndentedJSON(
			http.StatusInternalServerError,
			httpdto.GeneralResponse{
				Status: httpdto.StatusError,
			})
		return
	}

	tasks := make([]httpdto.TaskItem, 0, len(items))
	for _, item := range items {
		tasks = append(tasks, httpdto.TaskItem{
			ID:     *item.ItemID,
			Title:  *item.Title,
			IsDone: *item.IsDone,
		})
	}

	c.IndentedJSON(http.StatusOK, httpdto.GetTaskListResponse{
		GeneralResponse: httpdto.GeneralResponse{
			Status: httpdto.StatusOK,
		},
		Tasks: tasks,
	})

}

func (h *ToDoHandlers) HandlerGetTaskByID(c *gin.Context) {

	userID := c.GetUint64("userID")
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 0)

	if err != nil {
		handlers.SendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	item, err := h.itemGetter.GetByID(
		c.Request.Context(),
		coredto.User{
			UserID: &userID,
		},
		taskID,
	)

	if err != nil {
		handlers.SendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, httpdto.GetTaskByIDResponse{
		GeneralResponse: httpdto.GeneralResponse{
			Status: httpdto.StatusOK,
		},
		Task: httpdto.TaskItem{
			ID:     *item.ItemID,
			Title:  *item.Title,
			IsDone: *item.IsDone,
		},
	},
	)
}

func (h *ToDoHandlers) HandlerUpdateTaskByID(c *gin.Context) {

	userID := c.GetUint64("userID")
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 0)

	if err != nil {
		handlers.SendErrorResponse(c, http.StatusBadRequest)
		return
	}

	var changes httpdto.TaskItemChanges

	err = c.BindJSON(&changes)

	if err != nil {
		handlers.SendErrorResponse(c, http.StatusBadRequest)
		return
	}

	_, err = h.itemUpdater.Update(
		c.Request.Context(),
		coredto.ToDoItem{
			Owner: &coredto.User{
				UserID: &userID,
			},
			ItemID: &taskID,
			Title:  changes.Title,
			IsDone: changes.IsDone,
		})

	if err != nil {
		if errors.Is(err, todoprovider.ErrToDoNotFound) {
			handlers.SendErrorResponse(c, http.StatusNotFound)
			return
		}
		handlers.SendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, httpdto.GeneralResponse{
		Status: httpdto.StatusOK,
	})
}

func (h *ToDoHandlers) HandlerDeleteTaskByID(c *gin.Context) {

	userID := c.GetUint64("userID")
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 0)

	if err != nil {
		handlers.SendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	err = h.itemDeleter.Delete(
		c.Request.Context(),
		coredto.ToDoItem{
			Owner: &coredto.User{
				UserID: &userID,
			},
			ItemID: &taskID,
		},
	)

	if err != nil {
		if errors.Is(err, todoprovider.ErrToDoNotFound) {
			handlers.SendErrorResponse(c, http.StatusNotFound)
			return
		}
		handlers.SendErrorResponse(c, http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK,
		httpdto.GeneralResponse{
			Status: httpdto.StatusOK,
		},
	)
}
