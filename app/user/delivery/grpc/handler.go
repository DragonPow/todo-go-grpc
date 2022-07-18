package grpc

import (
	"context"
	"todo-go-grpc/app/adapter"
	"todo-go-grpc/app/user/domain"
	_usecase "todo-go-grpc/app/user/usecase"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	usecase _usecase.UserUsecase
	UnimplementedUserHandlerServer
}

func NewUserServerGrpc(gserver *grpc.Server, userUsecase _usecase.UserUsecase) {
	userServer := &server{
		usecase: userUsecase,
	}

	RegisterUserHandlerServer(gserver, userServer)
}

func (serverInstance *server) Login(ctx context.Context, req *LoginReq) (*BasicUser, error) {
	return nil, nil
}
func (serverInstance *server) Get(ctx context.Context, req *GetReq) (*User, error) {
	return nil, nil
}
func (serverInstance *server) Create(ctx context.Context, req *CreateReq) (*User, error) {
	new_user, err := serverInstance.usecase.Create(ctx, &domain.User{
		Name:     req.Name,
		Username: req.Username,
		Password: req.Password,
	})

	return *adapter.TransferDomainToProto(*new_user), err
}
func (serverInstance *server) Update(ctx context.Context, req *UpdateReq) (*User, error) {
	return nil, nil
}
func (serverInstance *server) Delete(ctx context.Context, req *DeleteReq) (*emptypb.Empty, error) {
	return nil, nil
}
