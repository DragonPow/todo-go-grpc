package internal

import (
	"context"
	"errors"
	"log"

	response_service "todo-go-grpc/app/response_handler"
	api "todo-go-grpc/app/user/api"
	domain "todo-go-grpc/app/user/domain"
	repository "todo-go-grpc/app/user/repository"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	repo repository.UserRepository
	api.UnimplementedUserHandlerServer
}

func RegisterGrpc(gserver *grpc.Server, repo repository.UserRepository) {
	userServer := &server{
		repo: repo,
	}

	api.RegisterUserHandlerServer(gserver, userServer)
}

func transferDomainToProto(in domain.User) *api.User {
	return &api.User{
		Id:          in.ID,
		Name:        in.Name,
		Username:    in.Username,
		Password:    in.Password,
		CreatedTime: timestamppb.New(in.CreatedAt),
	}
}

func transferProtoToDomain(in api.User) *domain.User {
	return &domain.User{
		ID:        in.Id,
		Name:      in.Name,
		Username:  in.Username,
		Password:  in.Password,
		CreatedAt: in.CreatedTime.AsTime(),
	}
}

func (serverInstance *server) Login(ctx context.Context, req *api.LoginReq) (*api.BasicUser, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	user, err := serverInstance.repo.GetByUsernameAndPassword(ctx, req.Username, req.Password)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrUserNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	// TODO: set token here

	user_basic := api.BasicUser{
		Id:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
	}

	return &user_basic, nil
}

func (serverInstance *server) Get(ctx context.Context, req *api.GetReq) (*api.User, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	user, err := serverInstance.repo.GetByID(ctx, req.Id)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrUserNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return transferDomainToProto(*user), nil
}

func (serverInstance *server) Create(ctx context.Context, req *api.CreateReq) (*api.User, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	new_user, err := serverInstance.repo.Create(ctx, &domain.User{
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

func (serverInstance *server) Update(ctx context.Context, req *api.UpdateReq) (*api.User, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	data := transferProtoToDomain(*req.NewUserInfor)
	new_user, err := serverInstance.repo.Update(ctx, req.Id, data)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrUserNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return transferDomainToProto(*new_user), nil
}

func (serverInstance *server) Delete(ctx context.Context, req *api.DeleteReq) (*emptypb.Empty, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	if err := serverInstance.repo.Delete(ctx, req.Id); err != nil {
		if errors.Is(err, domain.ErrUserNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return &emptypb.Empty{}, nil
}
