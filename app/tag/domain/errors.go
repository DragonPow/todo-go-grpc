package domain

import "errors"

var (
	ErrTagNotExists       = errors.New("ErrTagNotExists")
	ErrTagValueDuplicated = errors.New("ErrTagValueDuplicated")
	ErrTagStillReference  = errors.New("ErrTagStillReference")
)
