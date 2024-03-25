package userDTO

type UserDTO struct {
	Name    string `json:"name" validate:"required,correct_text"`
	Surname string `json:"surname" validate:"required,correct_text"`
	Mail    string `json:"mail" validate:"required,email"`
	Date    string `json:"date" validate:"required,date"`
	Cash    int    `json:"cash"`
}
