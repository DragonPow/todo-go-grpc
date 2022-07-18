package grpc

import (
	"context"
	"errors"
	"log"
	"todo-go-grpc/app/tag/domain"
	_usecase "todo-go-grpc/app/tag/usecase"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

func transferDomainToProto(in domain.Tag) *Tag {
	return &Tag{
		Id:          in.ID,
		Description: in.Description,
		CreatedTime: timestamppb.New(in.CreatedAt),
	}
}

func transferProtoToDomain(in Tag) *domain.Tag {
	return &domain.Tag{
		ID:          in.Id,
		Description: in.Description,
		CreatedAt:   in.CreatedTime.AsTime(),
	}
}

func (serverInstance *server) List(ctx context.Context, req *ListReq) (*ListTag, error) {
	tags_domain, err := serverInstance.usecase.FetchAll(ctx)

	if err != nil {
		log.Println(err.Error())
		return nil, grpc_status.Error(codes.Unknown, err.Error())
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
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
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
			return nil, grpc_status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return transferDomainToProto(*new_tag), nil
}

func (serverInstance *server) Update(ctx context.Context, req *UpdateReq) (*Tag, error) {
	data := transferProtoToDomain(*req.NewTagInfo)
	new_tag, err := serverInstance.usecase.Update(ctx, req.Id, data)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, domain.ErrTagIsExists) {
			return nil, grpc_status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return transferDomainToProto(*new_tag), nil
}

func (serverInstance *server) Delete(ctx context.Context, req *DeleteReq) (*emptypb.Empty, error) {
	err := serverInstance.usecase.Delete(ctx, req.Id)

	if err != nil {
		if errors.Is(err, domain.ErrTagNotExists) {
			return nil, grpc_status.Error(codes.NotFound, err.Error())
		}
		return nil, grpc_status.Error(codes.Unknown, err.Error())
	}

	return nil, nil
}
