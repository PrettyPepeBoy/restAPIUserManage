package userDTO

type UserDTO struct {
	Name    string `json:"name" validate:"required,correct_text"`
	Surname string `json:"surname" validate:"required,correct_text"`
	Mail    string `json:"mail" validate:"required,email"`
	Date    string `json:"date" validate:"required,date"`
	ID      int64  `json:"id,omitempty"`
	Cash    int    `json:"cash" validate:"required"`
}

type DTOUpdate struct {
	Name    *string `json:"name,omitempty"`
	Surname *string `json:"surname,omitempty"`
	Mail    *string `json:"mail,omitempty"`
	Date    *string `json:"date,omitempty"`
	Cash    *int    `json:"cash,omitempty"`
	ID      int64   `json:"id" validate:"required"`
}
