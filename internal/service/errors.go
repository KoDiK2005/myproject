package service

import "errors"

var (
	ErrForbidden  = errors.New("forbidden")
	ErrNotFound   = errors.New("not found")
	ErrEmailTaken = errors.New("email already taken")
)
