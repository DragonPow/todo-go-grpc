package usecase

import (
	"context"
	"fmt"
	"todo-go-grpc/app/dbservice"
	"todo-go-grpc/app/tag/domain"
	"todo-go-grpc/app/tag/repository"
)

type TagUsecase interface {
	FetchAll(ctx context.Context) ([]domain.Tag, error)
	GetByID(ctx context.Context, id int32) (*domain.Tag, error)
	Create(ctx context.Context, info *domain.Tag) (*domain.Tag, error)
	Update(ctx context.Context, id int32, new_info *domain.Tag) (*domain.Tag, error)
	Delete(ctx context.Context, id int32) error
	DeleteAll(ctx context.Context) error
}

type tagUsecase struct {
	db      dbservice.Database
	tagRepo repository.TagRepository
}

func NewTagUsecase(db dbservice.Database, t repository.TagRepository) TagUsecase {
	return &tagUsecase{
		db:      db,
		tagRepo: t,
	}
}

func (t *tagUsecase) FetchAll(ctx context.Context) ([]domain.Tag, error) {
	tags, err := t.tagRepo.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (t *tagUsecase) GetByID(ctx context.Context, id int32) (*domain.Tag, error) {
	tag, err := t.tagRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (t *tagUsecase) Create(ctx context.Context, info *domain.Tag) (*domain.Tag, error) {
	tag, err := t.tagRepo.Create(ctx, info)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (t *tagUsecase) Update(ctx context.Context, id int32, new_info *domain.Tag) (*domain.Tag, error) {
	tag, err := t.tagRepo.Update(ctx, id, new_info)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (t *tagUsecase) Delete(ctx context.Context, id int32) error {
	isSuccess := false
	tx := t.db.Db.Begin()

	defer func() {
		if isSuccess {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	// Check is exists
	if _, err := t.tagRepo.GetByID(ctx, id); err != nil {
		return err
	}

	if err := t.tagRepo.Delete(ctx, id); err != nil {
		return err
	}

	isSuccess = true
	return nil
}

func (t *tagUsecase) DeleteAll(ctx context.Context) error {
	return fmt.Errorf("Implemented needed")
}
