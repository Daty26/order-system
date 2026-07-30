package service

import "errors"

var (
	ErrInvalidStatus             = errors.New("invalid notification status")
	ErrInvalidMessage            = errors.New("message cannot be empty")
	ErrInvalidID                 = errors.New("invalid id input")
	ErrNotificationAlreadyExists = errors.New("notification already exists")
)
