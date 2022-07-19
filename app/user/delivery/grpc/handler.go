package grpc

import (
	"context"
	"errors"
	"log"
	"todo-go-grpc/app/user/domain"
	_usecase "todo-go-grpc/app/user/usecase"

	response_service "todo-go-grpc/app/responseservice"

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
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	user, err := serverInstance.usecase.Login(ctx, req.Username, req.Password)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrUsernameOrPasswordWrong) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
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
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	user, err := serverInstance.usecase.GetByID(ctx, req.Id)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrUserNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return transferDomainToProto(*user), nil
}

func (serverInstance *server) Create(ctx context.Context, req *CreateReq) (*User, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	new_user, err := serverInstance.usecase.Create(ctx, &domain.User{
		Name:     req.Name,
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrUserNameIsExists) {
			return nil, response_service.ResponseErrorInvalidArgument(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return transferDomainToProto(*new_user), nil
}

func (serverInstance *server) Update(ctx context.Context, req *UpdateReq) (*User, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	data := transferProtoToDomain(*req.NewUserInfor)
	new_user, err := serverInstance.usecase.Update(ctx, req.Id, data)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrUserNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return transferDomainToProto(*new_user), nil
}

func (serverInstance *server) Delete(ctx context.Context, req *DeleteReq) (*emptypb.Empty, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	if err := serverInstance.usecase.Delete(ctx, req.Id); err != nil {
		if errors.Is(err, domain.ErrUserNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return &emptypb.Empty{}, nil
}
