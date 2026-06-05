package repository

import "errors"

var (
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUser     = errors.New("duplicate user")
)
