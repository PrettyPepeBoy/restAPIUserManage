package service

import (
	"tstUser/internal/storage/storages"
)

type ProductService interface {
	CreateProducts(product storages.Product) (int64, error)
	UpdateProducts(product storages.Product) error
	GetProducts(id int64) (storages.Product, error)
	DeleteProduct(id int64) error
}

type UserService interface {
	CreateUser(user storages.User) (int64, error)
	DeleteUser(id int64) error
	GetUserInfo(Id int64) (storages.User, error)
	UpdateUser(user storages.User) error
}
