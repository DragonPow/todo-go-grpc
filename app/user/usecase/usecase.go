package usecase

import (
	"context"
	"todo-go-grpc/app/user/domain"
)

type UserUsecase interface {
	Login(ctx context.Context, username string, password string) (*domain.User, error)
	GetByID(ctx context.Context, id int32) (*domain.User, error)
	Create(ctx context.Context, info *domain.User) (*domain.User, error)
	Update(ctx context.Context, id int32, new_info *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id int32) error
}
