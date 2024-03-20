package userDTO

type UDTO struct {
	Name    string `json:"name" validate:"required,name"`
	Surname string `json:"surname" validate:"required,surname"`
	Mail    string `json:"mail" validate:"required,email"`
	Date    string `json:"date" validate:"required,date"`
	Cash    int    `json:"cash" validate:"required,cash"`
}
