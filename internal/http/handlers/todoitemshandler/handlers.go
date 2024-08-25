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

// HandlerCreateTask
// @Security 	ApiKeyAuth
// @Summary 	Create new task
// @Router 		/tasks [POST]
// @Param 		request body TaskItemChanges true "New task fields"
// @Tags 		TodoList
// @Produce		json
//
// @Success 200 		{object} 	GetTaskByIDResponse
// @Failure 400,401,500 {object}	GeneralResponse
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

// HandlerGetTaskList
// @Security 	ApiKeyAuth
// @Summary 	Get tasks list
// @Router 		/tasks [GET]
// @Tags 		TodoList
// @Produce		json
//
// @Success 200 	{object} 	GetTaskListResponse
// @Failure 401,500	{object}	GeneralResponse
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

// HandlerGetTaskByID
// @Security 	ApiKeyAuth
// @Summary 	Get single task by ID
// @Router 		/tasks/{id} [GET]
// @Param 		id	path int true "Task ID"
// @Tags 		TodoList
// @Produce		json
//
// @Success 200 			{object}	GetTaskByIDResponse
// @Failure 400,401,404,500 {object}	GeneralResponse
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

// HandlerUpdateTaskByID
// @Security 	ApiKeyAuth
// @Summary 	Change task fields by ID
// @Router 		/tasks/{id} [PATCH]
// @Param 		id	path int true "Task ID"
// @Param 		request body TaskItemChanges true "Fields changes"
// @Tags 		TodoList
// @Produce		json
//
// @Success 200 			{object}	GeneralResponse
// @Failure 400,401,404,500 {object}	GeneralResponse
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

// HandlerDeleteTaskByID
// @Security 	ApiKeyAuth
// @Summary 	Delete task by ID
// @Router 		/tasks/{id} [DELETE]
// @Param 		id	path int true "Task ID"
// @Tags 		TodoList
// @Produce		json
//
// @Success 200 			{object}	GeneralResponse
// @Failure 400,401,404,500 {object}	GeneralResponse
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
