package postgre

import (
	"context"
	"errors"
	"todo-go-grpc/app/dbservice"
	"todo-go-grpc/app/user/domain"
	"todo-go-grpc/app/user/repository"

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
	if err := u.Conn.Db.Create(&info).Error; err != nil {
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

	if err := u.Conn.Db.First(&new_info, id).Updates(&update).Error; err != nil {
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
