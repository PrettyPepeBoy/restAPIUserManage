package DTO

type ProductDTO struct {
	ID     int64  `json:"ID,omitempty"`
	Name   string `json:"name" validate:"required,correct_text"`
	Price  int    `json:"price" validate:"required"`
	Amount int    `json:"amount" validate:"required"`
}

type ProductDTOid struct {
	ID int64 `json:"id" validate:"required"`
}

type ProductDTOUpdate struct {
	ID     int64   `json:"id" validate:"required"`
	Name   *string `json:"name,omitempty"`
	Price  *int    `json:"price,omitempty"`
	Amount *int    `json:"amount,omitempty"`
}

type UserDTO struct {
	Name    string `json:"name" validate:"required,correct_text"`
	Surname string `json:"surname" validate:"required,correct_text"`
	Mail    string `json:"mail" validate:"required,email"`
	Date    string `json:"date" validate:"required,date"`
	ID      int64  `json:"id,omitempty"`
	Cash    int    `json:"cash" validate:"required"`
}

type UserDTOid struct {
	Id int64 `json:"id" validate:"required"`
}

type UserDTOUpdate struct {
	Name    *string `json:"name,omitempty"`
	Surname *string `json:"surname,omitempty"`
	Mail    *string `json:"mail,omitempty"`
	Date    *string `json:"date,omitempty"`
	Cash    *int    `json:"cash,omitempty"`
	ID      int64   `json:"id" validate:"required"`
}
