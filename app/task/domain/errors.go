package domain

import "errors"

var (
	ErrTaskNotExists = errors.New("ErrTaskNotExists")
	ErrTaskExists    = errors.New("ErrTaskExists")
	ErrTagNotExists  = errors.New("ErrTagNotExists")
)
