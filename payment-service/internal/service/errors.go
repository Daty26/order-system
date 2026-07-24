package service

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	// TODO prevent duplicate payments
	ErrPaymentAlreadyExists = errors.New("payment already exists")
)
