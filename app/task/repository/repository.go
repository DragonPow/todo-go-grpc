package repository

import (
	"context"
	"todo-go-grpc/app/task/domain"
)

type TaskRepository interface {
	Fetch(ctx context.Context, user_id int32, offset int32, number int32, conditions map[string]any) ([]domain.Task, error)
	GetByID(ctx context.Context, id int32) (*domain.Task, error)
	GetByUserId(ctx context.Context, user_id int32) ([]int32, error)
	IsExists(ctx context.Context, id int32) (bool, error)
	Create(ctx context.Context, user_id int32, info *domain.Task) (*domain.Task, error)
	Update(ctx context.Context, id int32, new_info *domain.Task, tags_add []int32, tags_remove []int32) (*domain.Task, error)
	Delete(ctx context.Context, ids []int32) error
}
