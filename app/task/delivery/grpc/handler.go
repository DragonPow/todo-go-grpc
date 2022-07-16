package grpc

import (
	usecase "todo-go-grpc/app/task/usecase"

	grpc "google.golang.org/grpc"
)

func NewTagServerGrpc(gserver *grpc.Server, tagUsecase usecase.TaskUsecase) {

}
