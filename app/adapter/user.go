package adapter

import (
	"todo-go-grpc/app/user/delivery/grpc"
	"todo-go-grpc/app/user/domain"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TransferDomainToProto(in domain.User) *grpc.User {
	return &grpc.User{
		Id:          in.ID,
		Name:        in.Name,
		Username:    in.Username,
		Password:    in.Password,
		CreatedTime: timestamppb.New(in.CreatedAt),
	}
}

func TransferProtoToDomain(in grpc.User) *domain.User {
	return &domain.User{
		ID:        in.Id,
		Name:      in.Name,
		Username:  in.Username,
		Password:  in.Password,
		CreatedAt: in.CreatedTime.AsTime(),
	}
}
