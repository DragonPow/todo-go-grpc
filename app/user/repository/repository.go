package repository

import (
	"context"
	"todo-go-grpc/app/user/domain"
)

type UserRepository interface {
	GetByUsernameAndPassword(ctx context.Context, username string, password string) (*domain.User, error)
	GetByID(ctx context.Context, id int32) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Create(ctx context.Context, info *domain.User) (*domain.User, error)
	Update(ctx context.Context, id int32, new_info *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id int32) error
}
