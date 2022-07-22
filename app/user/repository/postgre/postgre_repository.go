package postgre

import (
	"context"
	"errors"
	"todo-go-grpc/app/dbservice"
	"todo-go-grpc/app/user/domain"
	"todo-go-grpc/app/user/repository"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type userRepository struct {
	Conn dbservice.Database
}

func NewUserRepository(conn dbservice.Database) repository.UserRepository {
	return &userRepository{
		Conn: conn,
	}
}

func (u *userRepository) GetByUsernameAndPassword(ctx context.Context, username string, password string) (*domain.User, error) {
	var user domain.User
	if err := u.Conn.Db.Where("username=? AND password=?", username, password).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotExists
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (u *userRepository) GetByID(ctx context.Context, id int32) (*domain.User, error) {
	user := domain.User{ID: id}
	if err := u.Conn.Db.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotExists
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (u *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	if err := u.Conn.Db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotExists
		} else {
			return nil, err
		}
	}
	return &user, nil
}

func (u *userRepository) Create(ctx context.Context, info *domain.User) (*domain.User, error) {
	info_map := map[string]interface{}{
		"name":     info.Name,
		"username": info.Username,
		"password": info.Password,
	}
	if err := u.Conn.Db.Debug().Model(&domain.User{}).Create(info_map).Error; err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return nil, domain.ErrUserNameIsExists
			}
		}
		return nil, err
	}

	return info, nil
}

func (u *userRepository) Update(ctx context.Context, id int32, new_info *domain.User) (*domain.User, error) {
	update := map[string]any{
		"name":     new_info.Name,
		"username": new_info.Username,
		"password": new_info.Password,
	}

	if err := u.Conn.Db.Model(&new_info).Updates(&update).Error; err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return nil, domain.ErrUserNameIsExists
			}
		}
		return nil, err
	}

	return new_info, nil
}

func (u *userRepository) Delete(ctx context.Context, id int32) error {
	if err := u.Conn.Db.Delete(&domain.User{}, id).Error; err != nil {
		return err
	}

	return nil
}
