package storage

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExist    = errors.New("user already exists")
	ErrGoodNotFound = errors.New("goods not found")
	ErrGoodsEmpty   = errors.New("goods are empty")
	ErrGoodsExist   = errors.New("goods already exists")
)
