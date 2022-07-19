package usecase

import (
	"context"
	"todo-go-grpc/app/dbservice"
	tagRepo "todo-go-grpc/app/tag/repository"
	"todo-go-grpc/app/task/domain"
	taskRepo "todo-go-grpc/app/task/repository"
	userRepo "todo-go-grpc/app/user/repository"
)

type TaskUsecase interface {
	Fetch(ctx context.Context, user_id int32, offset int32, number int32, search_condition map[string]any) ([]domain.Task, error)
	GetByID(ctx context.Context, id int32) (*domain.Task, error)
	Create(ctx context.Context, user_id int32, info *domain.Task) (*domain.Task, error)
	Update(ctx context.Context, id int32, new_info *domain.Task, tags_add []int32, tags_remove []int32) (*domain.Task, error)
	Delete(ctx context.Context, ids []int32) error
	DeleteAllTaskOfUser(ctx context.Context, user_id int32) error
}

type taskUsecase struct {
	db       dbservice.Database
	taskRepo taskRepo.TaskRepository
	userRepo userRepo.UserRepository
	tagRepo  tagRepo.TagRepository
}

func NewTaskUsecase(db dbservice.Database, t taskRepo.TaskRepository, u userRepo.UserRepository, tag tagRepo.TagRepository) TaskUsecase {
	return &taskUsecase{
		db:       db,
		taskRepo: t,
		userRepo: u,
		tagRepo:  tag,
	}
}

func (t *taskUsecase) Fetch(ctx context.Context, user_id int32, offset int32, number int32, search_condition map[string]any) ([]domain.Task, error) {
	isSuccess := false
	tx := t.db.Db.Begin()

	defer func() {
		if isSuccess {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	// Check user exists
	if _, err := t.userRepo.GetByID(ctx, user_id); err != nil {
		return nil, err
	}

	// Fetch
	tasks, err := t.taskRepo.Fetch(ctx, user_id, offset, number, search_condition)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *taskUsecase) GetByID(ctx context.Context, id int32) (*domain.Task, error) {
	new_task, err := t.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return new_task, nil
}

func (t *taskUsecase) Create(ctx context.Context, creator_id int32, info *domain.Task) (*domain.Task, error) {
	isSuccess := false
	tx := t.db.Db.Begin()

	defer func() {
		if isSuccess {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	// Check user exists
	if _, err := t.userRepo.GetByID(ctx, creator_id); err != nil {
		return nil, err
	}

	task, err := t.taskRepo.Create(ctx, creator_id, info)
	if err != nil {
		return nil, err
	}

	isSuccess = true
	return task, nil
}

func (t *taskUsecase) Update(ctx context.Context, id int32, new_info *domain.Task, tags_add []int32, tags_remove []int32) (*domain.Task, error) {
	isSuccess := false
	tx := t.db.Db.Begin()

	defer func() {
		if isSuccess {
			tx.Commit()
		}
		tx.Rollback()
	}()

	// Check user exists
	if isExists, err := t.taskRepo.IsExists(ctx, id); err != nil {
		return nil, err
	} else if !isExists {
		return nil, domain.ErrTaskNotExists
	}

	return t.taskRepo.Update(ctx, id, new_info, tags_add, tags_remove)
}

func (t *taskUsecase) Delete(ctx context.Context, ids []int32) error {
	isSuccess := false
	tx := t.db.Db.Begin()

	defer func() {
		if isSuccess {
			tx.Commit()
		}
		tx.Rollback()
	}()

	// Check task exists
	// TODO: fix some logic wrong here
	// if isExists, err := t.taskRepo.isExists(ctx, ids); err != nil {
	// 	return err
	// }

	// Delete
	if err := t.taskRepo.Delete(ctx, ids); err != nil {
		return err
	}

	isSuccess = true
	return nil
}

func (t *taskUsecase) DeleteAllTaskOfUser(ctx context.Context, creator_id int32) error {
	isSuccess := false
	tx := t.db.Db.Begin()

	defer func() {
		if isSuccess {
			tx.Commit()
		}
		tx.Rollback()
	}()

	// Find id of user
	tasks, err := t.taskRepo.Fetch(ctx, creator_id, 0, -1, nil)
	if err != nil {
		return err
	}

	// Delete
	tasks_id := []int32{}
	for _, task := range tasks {
		tasks_id = append(tasks_id, task.ID)
	}
	if err := t.taskRepo.Delete(ctx, tasks_id); err != nil {
		return err
	}

	isSuccess = true
	return nil
}
