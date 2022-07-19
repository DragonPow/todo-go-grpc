package internal

import (
	"context"
	"errors"
	"log"
	response_service "todo-go-grpc/app/response_handler"
	api "todo-go-grpc/app/tag/api"
	domain "todo-go-grpc/app/tag/domain"
	repository "todo-go-grpc/app/tag/repository"

	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	repo repository.TagRepository
	api.UnimplementedTagHandlerServer
}

func RegisterGrpc(gserver *grpc.Server, repo repository.TagRepository) {
	tagServer := &server{
		repo: repo,
	}

	api.RegisterTagHandlerServer(gserver, tagServer)
}

func transferDomainToProto(in domain.Tag) *api.Tag {
	return &api.Tag{
		Id:          in.ID,
		Value:       in.Value,
		Description: in.Description,
		CreatedTime: timestamppb.New(in.CreatedAt),
	}
}

func transferProtoToDomain(in api.Tag) *domain.Tag {
	return &domain.Tag{
		ID:          in.Id,
		Description: in.Description,
		Value:       in.Value,
		CreatedAt:   in.CreatedTime.AsTime(),
	}
}

func (serverInstance *server) List(ctx context.Context, req *api.ListReq) (*api.ListTag, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	tags_domain, err := serverInstance.repo.FetchAll(ctx)

	if err != nil {
		log.Println(err.Error())
		return nil, response_service.ResponseErrorUnknown(err)
	}

	tags_rs := &api.ListTag{Tags: []*api.Tag{}}
	for _, task := range tags_domain {
		tags_rs.Tags = append(tags_rs.Tags, transferDomainToProto(task))
	}

	return tags_rs, nil
}

func (serverInstance *server) Get(ctx context.Context, req *api.GetReq) (*api.Tag, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	task, err := serverInstance.repo.GetByID(ctx, req.Id)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTagNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return transferDomainToProto(*task), nil
}

func (serverInstance *server) Create(ctx context.Context, req *api.CreateReq) (*api.Tag, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	new_tag, err := serverInstance.repo.Create(ctx, &domain.Tag{
		Value:       req.Value,
		Description: req.Description,
	})

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTagIsExists) {
			return nil, response_service.ResponseErrorAlreadyExists(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return transferDomainToProto(*new_tag), nil
}

func (serverInstance *server) Update(ctx context.Context, req *api.UpdateReq) (*api.Tag, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	data := transferProtoToDomain(*req.NewTagInfo)
	new_tag, err := serverInstance.repo.Update(ctx, req.Id, data)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTagIsExists) {
			return nil, response_service.ResponseErrorAlreadyExists(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return transferDomainToProto(*new_tag), nil
}

func (serverInstance *server) Delete(ctx context.Context, req *api.DeleteReq) (*emptypb.Empty, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	err := serverInstance.repo.Delete(ctx, req.Id)

	if err != nil {
		if errors.Is(err, domain.ErrTagNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return &emptypb.Empty{}, nil
}
