package grpc

import (
	"context"
	"errors"
	"log"
	"todo-go-grpc/app/tag/domain"
	_usecase "todo-go-grpc/app/tag/usecase"

	response_service "todo-go-grpc/app/responseservice"

	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	usecase _usecase.TagUsecase
	UnimplementedTagHandlerServer
}

func RegisterGrpc(gserver *grpc.Server, tagUsecase _usecase.TagUsecase) {
	tagServer := &server{
		usecase: tagUsecase,
	}

	RegisterTagHandlerServer(gserver, tagServer)
}

func transferDomainToProto(in domain.Tag) *Tag {
	return &Tag{
		Id:          in.ID,
		Value:       in.Value,
		Description: in.Description,
		CreatedTime: timestamppb.New(in.CreatedAt),
	}
}

func transferProtoToDomain(in Tag) *domain.Tag {
	return &domain.Tag{
		ID:          in.Id,
		Description: in.Description,
		Value:       in.Value,
		CreatedAt:   in.CreatedTime.AsTime(),
	}
}

func (serverInstance *server) List(ctx context.Context, req *ListReq) (*ListTag, error) {
	if err := req.Valid(); err != nil {
		return nil, response_service.ResponseErrorInvalidArgument(err)
	}

	tags_domain, err := serverInstance.usecase.FetchAll(ctx)

	if err != nil {
		log.Println(err.Error())
		return nil, response_service.ResponseErrorUnknown(err)
	}

	tags_rs := &ListTag{Tags: []*Tag{}}
	for _, task := range tags_domain {
		tags_rs.Tags = append(tags_rs.Tags, transferDomainToProto(task))
	}

	return tags_rs, nil
}

func (serverInstance *server) Get(ctx context.Context, req *GetReq) (*Tag, error) {
	task, err := serverInstance.usecase.GetByID(ctx, req.Id)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTagNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return transferDomainToProto(*task), nil
}

func (serverInstance *server) Create(ctx context.Context, req *CreateReq) (*Tag, error) {
	new_tag, err := serverInstance.usecase.Create(ctx, &domain.Tag{
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

func (serverInstance *server) Update(ctx context.Context, req *UpdateReq) (*Tag, error) {
	data := transferProtoToDomain(*req.NewTagInfo)
	new_tag, err := serverInstance.usecase.Update(ctx, req.Id, data)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTagIsExists) {
			return nil, response_service.ResponseErrorAlreadyExists(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return transferDomainToProto(*new_tag), nil
}

func (serverInstance *server) Delete(ctx context.Context, req *DeleteReq) (*emptypb.Empty, error) {
	err := serverInstance.usecase.Delete(ctx, req.Id)

	if err != nil {
		if errors.Is(err, domain.ErrTagNotExists) {
			return nil, response_service.ResponseErrorNotFound(err)
		}
		return nil, response_service.ResponseErrorUnknown(err)
	}

	return &emptypb.Empty{}, nil
}
