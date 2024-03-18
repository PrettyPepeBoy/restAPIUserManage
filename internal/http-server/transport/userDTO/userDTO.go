package userDTO

type UDTO struct {
	Name    string `json:"name" validate:"required,name"`
	Surname string `json:"surname" validate:"required,surname"`
	Cash    int    `json:"cash" validate:"required,cash"`
	Date    string `json:"date" validate:"required,date"`
}
