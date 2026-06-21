package service

import "errors"

var ErrInvalidInput = errors.New("The input is invalid")
var ErrInsufficientStock = errors.New("The stock amount is insufficient")