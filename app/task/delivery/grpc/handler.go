package grpc

import (
	"context"
	_usecase "todo-go-grpc/app/task/usecase"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	usecase _usecase.TaskUsecase
	UnimplementedTaskHandlerServer
}

func NewTagServerGrpc(gserver *grpc.Server, tagUsecase _usecase.TaskUsecase) {
	taskServer := &server{
		usecase: tagUsecase,
	}

	RegisterTaskHandlerServer(gserver, taskServer)
}

func (serverInstance *server) List(context.Context, *ListReq) (*ListTask, error) {
	return nil, nil
}

func (serverInstance *server) Get(context.Context, *GetReq) (*Task, error) {
	return nil, nil
}

func (serverInstance *server) Create(context.Context, *CreateReq) (*Task, error) {
	return nil, nil
}

func (serverInstance *server) Update(context.Context, *UpdateReq) (*Task, error) {
	return nil, nil
}

func (serverInstance *server) Delete(context.Context, *DeleteReq) (*emptypb.Empty, error) {
	return nil, nil
}

func (serverInstance *server) DeleteAll(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
