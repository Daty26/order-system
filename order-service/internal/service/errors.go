package service

import "errors"

var ErrInvalidRequest = errors.New("invalid order request")
var ErrProductNotFound = errors.New("product not found ")
var ErrForbiddenOrder = errors.New("forbidden order")
var ErrOrderCannotBeCanceled = errors.New("cannot cancel the order")
