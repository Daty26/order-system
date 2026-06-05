package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrIncorrectID        = errors.New("user with such id doesn't exist")
	ErrInvalidUserInput   = errors.New("invalid user input")
	ErrNotFound           = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
)
