package mocks

import (
	"context"
	"fmt"
	todoprotobufv1 "github.com/IldarGaleev/todo-backend-service/pkg/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type ToDoGrpcMock struct {
	isDown    bool
	blackList map[string]struct{}
}

func New(isDown bool) *ToDoGrpcMock {
	return &ToDoGrpcMock{
		isDown:    isDown,
		blackList: make(map[string]struct{}, 2),
	}
}

func decodeToken(token string) ([]string, error) {
	tokenData := strings.Split(token, ":")

	if tokenData == nil || len(tokenData) != 2 {
		return nil, fmt.Errorf("invalid token")
	}
	return tokenData, nil
}

func (p ToDoGrpcMock) Login(
	ctx context.Context,
	in *todoprotobufv1.LoginRequest,
	opts ...grpc.CallOption) (*todoprotobufv1.LoginResponce, error) {
	if p.isDown {
		return nil, fmt.Errorf("service is down")
	}

	if in.GetEmail() == "user1" && in.GetPassword() == "pass" {
		return &todoprotobufv1.LoginResponce{
			Token: "1:user1",
		}, nil

	}
	return nil, status.Error(codes.PermissionDenied, "Username or password incorrect")
}

func (p ToDoGrpcMock) Logout(
	ctx context.Context,
	in *todoprotobufv1.LogoutRequest,
	opts ...grpc.CallOption) (*todoprotobufv1.LogoutResponce, error) {
	if p.isDown {
		return nil, fmt.Errorf("service is down")
	}
	tokenData, err := decodeToken(in.GetToken())

	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "Bad token")
	}

	if _, ok := p.blackList[tokenData[0]]; !ok && tokenData[0] == "1" && tokenData[1] == "user1" {
		p.blackList[tokenData[0]] = struct{}{}
		return &todoprotobufv1.LogoutResponce{
			Success: true,
		}, nil
	}

	return &todoprotobufv1.LogoutResponce{
		Success: false,
	}, status.Error(codes.Unauthenticated, "User not found")

}

func (p ToDoGrpcMock) CheckSecret(
	ctx context.Context,
	in *todoprotobufv1.CheckSecretRequest,
	opts ...grpc.CallOption) (*todoprotobufv1.CheckSecretResponce, error) {
	if p.isDown {
		return nil, fmt.Errorf("service is down")
	}

	tokenData, err := decodeToken(in.GetSecret())
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "Bad token")
	}

	if tokenData[0] == "1" {
		if _, ok := p.blackList[tokenData[0]]; !ok && tokenData[1] == "user1" {
			return &todoprotobufv1.CheckSecretResponce{
				Email:  "user1",
				UserId: 1,
			}, nil
		}
	}

	return nil, status.Error(codes.Unauthenticated, "Permission denied")
}

func (p ToDoGrpcMock) CreateTask(
	ctx context.Context,
	in *todoprotobufv1.CreateTaskRequest,
	opts ...grpc.CallOption) (*todoprotobufv1.CreateTaskResponce, error) {
	if p.isDown {
		return nil, fmt.Errorf("service is down")
	}
	//TODO implement me
	panic("implement me")
}

func (p ToDoGrpcMock) ListTasks(
	ctx context.Context,
	in *todoprotobufv1.ListTasksRequest,
	opts ...grpc.CallOption) (*todoprotobufv1.ListTasksResponce, error) {
	if p.isDown {
		return nil, fmt.Errorf("service is down")
	}
	//TODO implement me
	panic("implement me")
}

func (p ToDoGrpcMock) GetTaskByID(
	ctx context.Context,
	in *todoprotobufv1.TaskByIdRequest,
	opts ...grpc.CallOption) (*todoprotobufv1.GetTaskByIdResponce, error) {
	if p.isDown {
		return nil, fmt.Errorf("service is down")
	}
	//TODO implement me
	panic("implement me")
}

func (p ToDoGrpcMock) UpdateTaskByID(
	ctx context.Context,
	in *todoprotobufv1.UpdateTaskByIdRequest,
	opts ...grpc.CallOption) (*todoprotobufv1.ChangedTaskByIdResponce, error) {
	if p.isDown {
		return nil, fmt.Errorf("service is down")
	}
	//TODO implement me
	panic("implement me")
}

func (p ToDoGrpcMock) DeleteTaskByID(
	ctx context.Context,
	in *todoprotobufv1.TaskByIdRequest,
	opts ...grpc.CallOption) (*todoprotobufv1.ChangedTaskByIdResponce, error) {
	if p.isDown {
		return nil, fmt.Errorf("service is down")
	}
	//TODO implement me
	panic("implement me")
}
