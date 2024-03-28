package storage

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExist       = errors.New("user already exists")
	ErrProductNotFound = errors.New("products not found")
	ErrProductsEmpty   = errors.New("products are empty")
	ErrProductsExist   = errors.New("products already exists")
	ErrCashIsNotEnough = errors.New("cash is not enough")
)
