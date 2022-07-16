package usecase

import (
	"context"
	"todo-go-grpc/app/task/domain"
)

type TaskUsecase interface {
	Fetch(ctx context.Context, user_id int32, start_index int32, number int32, search_condition map[string]any) ([]domain.Task, error)
	GetByID(ctx context.Context, id int32) (*domain.Task, error)
	Create(ctx context.Context, user_id int32, info *domain.Task) (*domain.Task, error)
	Update(ctx context.Context, id int32, new_info *domain.Task) (*domain.Task, error)
	Delete(ctx context.Context, ids []int32) error
	DeleteAllTaskOfUser(ctx context.Context, user_id int32) error
}
