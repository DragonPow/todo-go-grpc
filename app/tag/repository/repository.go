package repository

import (
	"context"
	"todo-go-grpc/app/tag/domain"
)

type TagRepository interface {
	FetchAll(ctx context.Context) ([]domain.Tag, error)
	GetByID(ctx context.Context, id int32) (*domain.Tag, error)
	Create(ctx context.Context, info *domain.Tag) (*domain.Tag, error)
	Update(ctx context.Context, id int32, new_info *domain.Tag) (*domain.Tag, error)
	Delete(ctx context.Context, id int32) error
	DeleteAll(ctx context.Context) error
}
