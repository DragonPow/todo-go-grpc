package domain

import "errors"

var (
	ErrUserNotExists           = errors.New("ErrUserNotExists")
	ErrUserNameIsExists        = errors.New("ErrUserIsExists")
	ErrUsernameOrPasswordWrong = errors.New("ErrUsernameOrPasswordWrong")
)
