package postgre

import (
	"context"
	"errors"
	"fmt"

	"todo-go-grpc/app/dbservice"
	"todo-go-grpc/app/task/domain"
	"todo-go-grpc/app/task/repository"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type tagRepository struct {
	Conn dbservice.Database
}

func NewTagRepository(conn dbservice.Database) repository.TagRepository {
	return &tagRepository{
		Conn: conn,
	}
}

func (t *tagRepository) FetchAll(ctx context.Context) ([]domain.Tag, error) {
	var tags []domain.Tag
	if err := t.Conn.Db.Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (t *tagRepository) GetByID(ctx context.Context, id int32) (*domain.Tag, error) {
	var tag domain.Tag
	if err := t.Conn.Db.First(&tag, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTagNotExists
		}
		return nil, err
	}

	return &tag, nil
}

func (t *tagRepository) Create(ctx context.Context, info *domain.Tag) (*domain.Tag, error) {
	if err := t.Conn.Db.Create(&info).Error; err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok && errors.Is(err, pgError) {
			// Duplicate value
			if pgError.Code == "23505" {
				return nil, domain.ErrTagIsExists
			}
		}
		return nil, err
	}

	return info, nil
}

func (t *tagRepository) Update(ctx context.Context, id int32, new_info *domain.Tag) (*domain.Tag, error) {
	new_info_map := map[string]any{}
	new_info_map["value"] = new_info.Value
	new_info_map["description"] = new_info.Description

	if err := t.Conn.Db.First(&new_info, id).Updates(new_info_map).Error; err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok && errors.Is(err, pgError) {
			// Value of tag is duplicate
			if pgError.Code == "23505" {
				return nil, domain.ErrTagIsExists
			}
		}
		return nil, err
	}

	return new_info, nil
}

func (t *tagRepository) Delete(ctx context.Context, id int32) error {
	tag := domain.Tag{ID: id}
	if err := t.Conn.Db.Delete(&tag).Error; err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok && errors.Is(err, pgError) {
			// Tag still another reference
			if pgError.Code == "23503" {
				return domain.ErrTagStillReference
			}
		}
		return err
	}

	return nil
}

func (t *tagRepository) DeleteAll(ctx context.Context) error {
	return fmt.Errorf("Implemented needed")
}
