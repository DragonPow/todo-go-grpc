package domain

import "errors"

var (
	// Task
	ErrTaskNotExists = errors.New("ErrTaskNotExists")
	ErrTaskExists    = errors.New("ErrTaskExists")
	ErrUserNotExists = errors.New("ErrUserNotExists")

	// Tag
	ErrTagNotExists      = errors.New("ErrTagNotExists")
	ErrTagIsExists       = errors.New("ErrTagIsExists")
	ErrTagStillReference = errors.New("ErrTagStillReference")
)
