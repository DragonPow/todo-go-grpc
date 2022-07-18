package grpc

import (
	"context"
	_usecase "todo-go-grpc/app/tag/usecase"

	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	usecase _usecase.TagUsecase
	UnimplementedTagHandlerServer
}

func NewTagServerGrpc(gserver *grpc.Server, tagUsecase _usecase.TagUsecase) {
	tagServer := &server{
		usecase: tagUsecase,
	}

	RegisterTagHandlerServer(gserver, tagServer)
}

func (serverInstance *server) List(ctx context.Context, request *ListReq) (*ListTag, error) {
	return nil, nil
}

func (serverInstance *server) Get(ctx context.Context, request *GetReq) (*Tag, error) {
	return nil, nil
}

func (serverInstance *server) Create(ctx context.Context, request *CreateReq) (*Tag, error) {
	return nil, nil
}

func (serverInstance *server) Update(ctx context.Context, request *UpdateReq) (*Tag, error) {
	return nil, nil
}

func (serverInstance *server) Delete(ctx context.Context, request *DeleteReq) (*emptypb.Empty, error) {
	return nil, nil
}
