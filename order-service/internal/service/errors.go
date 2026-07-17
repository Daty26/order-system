package service

import "errors"

var ErrInvalidOrder = errors.New("invalid order")
var ErrProductNotFound = errors.New("product not found ")
var ErrForbiddenOrder = errors.New("forbidden order")
var ErrOrderCannotBeCanceled = errors.New("cannot cancel the order")
