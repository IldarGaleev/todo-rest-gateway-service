package main

import (
	"context"
	"net/http"
	"strconv"

	todo_protobuf_v1 "github.com/IldarGaleev/todo-backend-service/pkg/grpc/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GeneralResponseStatus string

const (
	StatusOK    = GeneralResponseStatus("ok")
	StatusError = GeneralResponseStatus("error")
)

type GeneralResponse struct {
	Status GeneralResponseStatus `json:"status"`
}

type MyIDResponse struct {
	GeneralResponse
	ID int `json:"id"`
}

type TaskItem struct {
	ID     uint64 `json:"id"`
	Title  string `json:"title"`
	IsDone bool   `json:"is_done"`
}

type TaskItemChanges struct {
	Title  *string `json:"title,omitempty"`
	IsDone *bool   `json:"is_done,omitempty"`
}

type GetTaskListResponse struct {
	GeneralResponse
	Tasks []TaskItem `json:"tasks"`
}

type GetTaskByIDResponse struct {
	GeneralResponse
	Task TaskItem `json:"task"`
}

func SendErrorResponse(c *gin.Context, code int) {
	c.IndentedJSON(
		code,
		GeneralResponse{
			Status: StatusError,
		})
}

func CreateHandlerGetTaskList(client todo_protobuf_v1.ToDoServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := client.ListTasks(c.Request.Context(), &todo_protobuf_v1.ListTasksRequest{
			UserId: 1,
		})

		if err != nil {
			c.IndentedJSON(
				http.StatusInternalServerError,
				GeneralResponse{
					Status: StatusError,
				})
			return
		}

		rTasks := resp.GetTasks()
		tasks := make([]TaskItem, 0, len(rTasks))
		for _, item := range rTasks {
			tasks = append(tasks, TaskItem{
				ID:     item.TaskId,
				Title:  item.Title,
				IsDone: item.IsDone,
			})
		}

		c.IndentedJSON(http.StatusOK, GetTaskListResponse{
			GeneralResponse: GeneralResponse{
				Status: StatusOK,
			},
			Tasks: tasks,
		})
	}
}

func CreateHandlerGetTaskByID(client todo_protobuf_v1.ToDoServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		//token := c.Request.Header.Get("Authorization")
		//TODO: check auth

		taskID, err := strconv.ParseUint(c.Param("id"), 10, 0)

		if err != nil {
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		resp, err := client.GetTaskByID(c.Request.Context(), &todo_protobuf_v1.TaskByIdRequest{
			TaskId: taskID,
			UserId: 1,
		})

		if err != nil {
			if status.Code(err) == codes.NotFound {
				SendErrorResponse(c, http.StatusNotFound)
				return
			}
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		c.IndentedJSON(http.StatusOK, GetTaskByIDResponse{
			GeneralResponse: GeneralResponse{
				Status: StatusOK,
			},
			Task: TaskItem{
				ID:     resp.GetTaskId(),
				Title:  resp.GetTitle(),
				IsDone: resp.GetIsDone(),
			},
		},
		)
	}
}

func CreateHandlerDeleteTaskByID(client todo_protobuf_v1.ToDoServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		//token := c.Request.Header.Get("Authorization")
		//TODO: check auth

		taskID, err := strconv.ParseUint(c.Param("id"), 10, 0)

		if err != nil {
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		_, err = client.DeleteTaskByID(c.Request.Context(), &todo_protobuf_v1.TaskByIdRequest{
			TaskId: taskID,
			UserId: 1,
		})

		if err != nil {
			if status.Code(err) == codes.NotFound {
				SendErrorResponse(c, http.StatusNotFound)
				return
			}
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		c.IndentedJSON(http.StatusOK,
			GeneralResponse{
				Status: StatusOK,
			},
		)
	}
}

func CreateHandlerUpdateTaskByID(client todo_protobuf_v1.ToDoServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		//token := c.Request.Header.Get("Authorization")
		//TODO: check auth

		taskID, err := strconv.ParseUint(c.Param("id"), 10, 0)
		var changes TaskItemChanges
		c.BindJSON(&changes)

		if err != nil {
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		_, err = client.UpdateTaskByID(c.Request.Context(), &todo_protobuf_v1.UpdateTaskByIdRequest{
			TaskId: taskID,
			UserId: 1,
			Title:  changes.Title,
			IsDone: changes.IsDone,
		})

		if err != nil {
			if status.Code(err) == codes.NotFound {
				SendErrorResponse(c, http.StatusNotFound)
				return
			}
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		c.IndentedJSON(http.StatusOK, GeneralResponse{
			Status: StatusOK,
		})
	}
}

func CreateHandlerCreateTask(client todo_protobuf_v1.ToDoServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		//token := c.Request.Header.Get("Authorization")
		//TODO: check auth

		var changes TaskItemChanges
		c.BindJSON(&changes)


		resp, err := client.CreateTask(c.Request.Context(), &todo_protobuf_v1.CreateTaskRequest{
			Title: *changes.Title,
			UserId: 1,
		})

		if err != nil {
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		c.IndentedJSON(http.StatusOK,GetTaskByIDResponse{
			GeneralResponse: GeneralResponse{
				Status: StatusOK,
			},
			Task: TaskItem{
				ID: resp.GetTaskId(),
				Title: *changes.Title,
			},
		})
	}
}

func unaryInterceptor(
	ctx context.Context,
	method string, req, reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {

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

func main() {
	router := gin.Default()

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithUnaryInterceptor(unaryInterceptor))

	conn, err := grpc.NewClient("localhost:9090", opts...)

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := todo_protobuf_v1.NewToDoServiceClient(conn)

	router.POST("/api/v1/tasks", CreateHandlerCreateTask(client))
	router.GET("/api/v1/tasks", CreateHandlerGetTaskList(client))
	router.GET("/api/v1/tasks/:id", CreateHandlerGetTaskByID(client))
	router.PATCH("/api/v1/tasks/:id", CreateHandlerUpdateTaskByID(client))
	router.DELETE("/api/v1/tasks/:id", CreateHandlerDeleteTaskByID(client))

	router.Run("localhost:8080")
}
