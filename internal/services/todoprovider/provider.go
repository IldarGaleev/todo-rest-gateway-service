// Package todoprovider implements ToDo list items gRPC bindings
package todoprovider

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
	ErrToDoInternal = errors.New("todo provider internal error")
	ErrToDoNotFound = errors.New("todo item not found")
)

type ToDoProvider struct {
	logger *slog.Logger
	client todoprotobufv1.ToDoServiceClient
}

func New(
	logger *slog.Logger,
	client todoprotobufv1.ToDoServiceClient,
) *ToDoProvider {
	return &ToDoProvider{
		logger: logger.With("module", "todoprovider"),
		client: client,
	}
}

func (p *ToDoProvider) Create(
	ctx context.Context,
	owner coredto.User,
	title string,
) (*coredto.ToDoItem, error) {
	log := p.logger.With("method", "Create")
	resp, err := p.client.CreateTask(
		ctx,
		&todoprotobufv1.CreateTaskRequest{
			Title:  title,
			UserId: *owner.UserID,
		},
	)

	if err != nil {
		log.Error("create error", slog.Any("err", err))
		return nil, errors.Join(ErrToDoInternal, err)
	}

	newItemID := resp.GetTaskId()
	isDone := false

	return &coredto.ToDoItem{
		ItemID: &newItemID,
		Owner:  &owner,
		Title:  &title,
		IsDone: &isDone,
	}, nil
}

func (p *ToDoProvider) Delete(
	ctx context.Context,
	item coredto.ToDoItem,
) error {
	log := p.logger.With("method", "Delete")
	_, err := p.client.DeleteTaskByID(ctx, &todoprotobufv1.TaskByIdRequest{
		TaskId: *item.ItemID,
		UserId: *item.Owner.UserID,
	})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return ErrToDoNotFound
		}
		log.Error("delete error", slog.Any("err", err))
		return errors.Join(ErrToDoInternal, err)
	}
	return nil
}

func (p *ToDoProvider) GetByID(
	ctx context.Context,
	owner coredto.User,
	itemID uint64,
) (*coredto.ToDoItem, error) {
	log := p.logger.With("method", "GetByID")
	resp, err := p.client.GetTaskByID(ctx, &todoprotobufv1.TaskByIdRequest{
		TaskId: itemID,
		UserId: *owner.UserID,
	})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrToDoNotFound
		}
		log.Error("get by id error", slog.Any("err", err))
		return nil, errors.Join(ErrToDoInternal, err)
	}

	title := resp.GetTitle()
	isDone := resp.GetIsDone()

	return &coredto.ToDoItem{
		ItemID: &itemID,
		Title:  &title,
		IsDone: &isDone,
		Owner:  &owner,
	}, nil

}

func (p *ToDoProvider) GetList(
	ctx context.Context,
	owner coredto.User,
) ([]coredto.ToDoItem, error) {
	log := p.logger.With("method", "GetList")
	resp, err := p.client.ListTasks(ctx,
		&todoprotobufv1.ListTasksRequest{
			UserId: *owner.UserID,
		})

	if err != nil {
		log.Error("get list error", slog.Any("err", err))
		return nil, errors.Join(ErrToDoInternal, err)
	}

	tasksR := resp.GetTasks()
	result := make([]coredto.ToDoItem, 0, len(tasksR))

	for _, item := range tasksR {
		result = append(result, coredto.ToDoItem{
			ItemID: &item.TaskId,
			Title:  &item.Title,
			IsDone: &item.IsDone,
			Owner:  &owner,
		})
	}

	return result, nil
}

func (p *ToDoProvider) Update(
	ctx context.Context,
	item coredto.ToDoItem,
) (*coredto.ToDoItem, error) {
	log := p.logger.With("method", "Update")
	_, err := p.client.UpdateTaskByID(
		ctx,
		&todoprotobufv1.UpdateTaskByIdRequest{
			TaskId: *item.ItemID,
			UserId: *item.Owner.UserID,
			Title:  item.Title,
			IsDone: item.IsDone,
		})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrToDoNotFound
		}
		log.Error("update error", slog.Any("err", err))
		return nil, errors.Join(ErrToDoInternal, err)
	}

	//TODO:update gRPC for get updated item

	return &item, nil
}
