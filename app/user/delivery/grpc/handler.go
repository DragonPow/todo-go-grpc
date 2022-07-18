package grpc

import (
	"context"
	"log"
	"todo-go-grpc/app/user/domain"
	_usecase "todo-go-grpc/app/user/usecase"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	usecase _usecase.UserUsecase
	UnimplementedUserHandlerServer
}

func RegisterGrpc(gserver *grpc.Server, userUsecase _usecase.UserUsecase) {
	userServer := &server{
		usecase: userUsecase,
	}

	RegisterUserHandlerServer(gserver, userServer)
}

func transferDomainToProto(in domain.User) *User {
	return &User{
		Id:          in.ID,
		Name:        in.Name,
		Username:    in.Username,
		Password:    in.Password,
		CreatedTime: timestamppb.New(in.CreatedAt),
	}
}

func transferProtoToDomain(in User) *domain.User {
	return &domain.User{
		ID:        in.Id,
		Name:      in.Name,
		Username:  in.Username,
		Password:  in.Password,
		CreatedAt: in.CreatedTime.AsTime(),
	}
}

func (serverInstance *server) Login(ctx context.Context, req *LoginReq) (*BasicUser, error) {
	user, err := serverInstance.usecase.Login(ctx, req.Username, req.Password)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// TODO: set token here

	user_basic := BasicUser{
		Id:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
	}

	return &user_basic, nil
}
func (serverInstance *server) Get(ctx context.Context, req *GetReq) (*User, error) {
	user, err := serverInstance.usecase.GetByID(ctx, req.Id)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return transferDomainToProto(*user), nil
}
func (serverInstance *server) Create(ctx context.Context, req *CreateReq) (*User, error) {
	new_user, err := serverInstance.usecase.Create(ctx, &domain.User{
		Name:     req.Name,
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return transferDomainToProto(*new_user), nil
}
func (serverInstance *server) Update(ctx context.Context, req *UpdateReq) (*User, error) {
	data := transferProtoToDomain(*req.NewUserInfor)
	new_user, err := serverInstance.usecase.Update(ctx, req.Id, data)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return transferDomainToProto(*new_user), nil
}
func (serverInstance *server) Delete(ctx context.Context, req *DeleteReq) (*emptypb.Empty, error) {
	err := serverInstance.usecase.Delete(ctx, req.Id)

	return nil, err
}
