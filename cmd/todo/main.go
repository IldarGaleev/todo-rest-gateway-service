package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

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

type LoginResponse struct {
	GeneralResponse
	Token string `json:"token"`
	// RefreshToken string `json:"refresh_token"`
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
			UserId: c.GetUint64("userID"),
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

		taskID, err := strconv.ParseUint(c.Param("id"), 10, 0)

		if err != nil {
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		resp, err := client.GetTaskByID(c.Request.Context(), &todo_protobuf_v1.TaskByIdRequest{
			TaskId: taskID,
			UserId: c.GetUint64("userID"),
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

		taskID, err := strconv.ParseUint(c.Param("id"), 10, 0)

		if err != nil {
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		_, err = client.DeleteTaskByID(c.Request.Context(), &todo_protobuf_v1.TaskByIdRequest{
			TaskId: taskID,
			UserId: c.GetUint64("userID"),
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

		taskID, err := strconv.ParseUint(c.Param("id"), 10, 0)
		var changes TaskItemChanges
		c.BindJSON(&changes)

		if err != nil {
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		_, err = client.UpdateTaskByID(c.Request.Context(), &todo_protobuf_v1.UpdateTaskByIdRequest{
			TaskId: taskID,
			UserId: c.GetUint64("userID"),
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
		var changes TaskItemChanges
		c.BindJSON(&changes)

		resp, err := client.CreateTask(c.Request.Context(), &todo_protobuf_v1.CreateTaskRequest{
			Title:  *changes.Title,
			UserId: c.GetUint64("userID"),
		})

		if err != nil {
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		c.IndentedJSON(http.StatusOK, GetTaskByIDResponse{
			GeneralResponse: GeneralResponse{
				Status: StatusOK,
			},
			Task: TaskItem{
				ID:    resp.GetTaskId(),
				Title: *changes.Title,
			},
		})
	}
}

func CreateJWTAuthMiddleware(client todo_protobuf_v1.ToDoServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		
		if token == "" {
			SendErrorResponse(c,http.StatusUnauthorized)
			c.Abort()
			return
		}

		parts:=strings.Split(token, " ")
		method, token := parts[0], parts[1]

		if strings.ToLower(method) != "bearer"{
			SendErrorResponse(c,http.StatusBadRequest)
			c.Abort()
			return
		}

		resp, err := client.CheckSecret(c.Request.Context(), &todo_protobuf_v1.CheckSecretRequest{
			Secret: token,
		})

		if err != nil {
			SendErrorResponse(c,http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("userID",resp.UserId)
		c.Set("jwtToken",token)
		c.Next()
	}
}

func CreateHandlerLogin(client todo_protobuf_v1.ToDoServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		email, pass, hasAuth := c.Request.BasicAuth()

		if !hasAuth {
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			SendErrorResponse(c, http.StatusUnauthorized)
			return
		}

		resp, err := client.Login(c.Request.Context(), &todo_protobuf_v1.LoginRequest{
			Email:    email,
			Password: pass,
		})

		if err != nil {
			if status.Code(err) == codes.PermissionDenied {
				SendErrorResponse(c, http.StatusUnauthorized)
				return
			}
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		c.IndentedJSON(http.StatusOK, LoginResponse{
			GeneralResponse: GeneralResponse{
				Status: StatusOK,
			},
			Token: resp.GetToken(),
		})
	}
}

func CreateHandlerLogout(client todo_protobuf_v1.ToDoServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		_, err := client.Logout(c.Request.Context(), &todo_protobuf_v1.LogoutRequest{
			Token: c.GetString("jwtToken"),
		})

		if err != nil {
			SendErrorResponse(c, http.StatusInternalServerError)
			return
		}

		c.IndentedJSON(http.StatusOK, GeneralResponse{
			Status: StatusOK,
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

	const apiBasePath = "/api/v1/"
	
	apiAuth := router.Group(apiBasePath)
	apiAuth.Use(CreateJWTAuthMiddleware(client))
	{
		apiAuth.POST("/tasks", CreateHandlerCreateTask(client))
		apiAuth.GET("/tasks", CreateHandlerGetTaskList(client))
		apiAuth.GET("/tasks/:id", CreateHandlerGetTaskByID(client))
		apiAuth.PATCH("/tasks/:id", CreateHandlerUpdateTaskByID(client))
		apiAuth.DELETE("/tasks/:id", CreateHandlerDeleteTaskByID(client))
		apiAuth.GET("/logout", CreateHandlerLogout(client))
	}

	apiNoAuth := router.Group(apiBasePath)
	{
		apiNoAuth.POST("/login", CreateHandlerLogin(client))
	}

	router.Run("localhost:8080")
}
