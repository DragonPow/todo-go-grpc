package usecase

import (
	"context"
	"errors"
	"todo-go-grpc/app/dbservice"
	"todo-go-grpc/app/user/domain"
	"todo-go-grpc/app/user/repository"
)

type UserUsecase interface {
	Login(ctx context.Context, username string, password string) (*domain.User, error)
	GetByID(ctx context.Context, id int32) (*domain.User, error)
	Create(ctx context.Context, info *domain.User) (*domain.User, error)
	Update(ctx context.Context, id int32, new_info *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id int32) error
}

type userUsecase struct {
	db       dbservice.Database
	userRepo repository.UserRepository
}

func NewUserUsecase(db dbservice.Database, u repository.UserRepository) UserUsecase {
	return &userUsecase{
		db:       db,
		userRepo: u,
	}
}

func (u *userUsecase) Login(ctx context.Context, username string, password string) (*domain.User, error) {
	user, err := u.userRepo.GetByUsernameAndPassword(ctx, username, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Create(ctx context.Context, info *domain.User) (*domain.User, error) {
	isSuccess := false
	tx := u.db.Db.Begin()

	defer func() {
		if isSuccess {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if _, err := u.userRepo.GetByUsername(ctx, info.Username); err != nil {
		// If username not found, add it
		if errors.Is(err, domain.ErrUserNotExists) {
			new_user, new_user_err := u.userRepo.Create(ctx, info)
			if new_user_err != nil {
				return nil, new_user_err
			}

			isSuccess = true
			return new_user, nil
		} else {
			return nil, err
		}
	} else {
		return nil, domain.ErrUserNameIsExists
	}
}

func (u *userUsecase) Update(ctx context.Context, id int32, new_info *domain.User) (*domain.User, error) {
	isSuccess := false
	tx := u.db.Db.Begin()

	defer func() {
		if isSuccess {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	// Check user is exists
	if _, err := u.userRepo.GetByID(ctx, id); err != nil {
		return nil, err
	}

	isSuccess = true
	return u.userRepo.Update(ctx, id, new_info)
}

func (u *userUsecase) Delete(ctx context.Context, id int32) error {
	isSuccess := false
	tx := u.db.Db.Begin()

	defer func() {
		if isSuccess {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	// Check ID is exists
	if _, err := u.userRepo.GetByID(ctx, id); err != nil {
		return err
	}

	// Delete
	if err := u.userRepo.Delete(ctx, id); err != nil {
		return err
	} else {
		return nil
	}
}

func (u *userUsecase) GetByID(ctx context.Context, id int32) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
