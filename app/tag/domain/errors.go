package domain

import "errors"

var (
	ErrTagNotExists      = errors.New("ErrTagNotExists")
	ErrTagIsExists       = errors.New("ErrTagIsExists")
	ErrTagStillReference = errors.New("ErrTagStillReference")
)
