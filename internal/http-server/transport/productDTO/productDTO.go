package productDTO

type DTOWithID struct {
	ID     int64  `json:"ID,omitempty"`
	Name   string `json:"name" validate:"required,correct_text"`
	Price  int64  `json:"price" validate:"required"`
	Amount int64  `json:"amount" validate:"required"`
}

type DTOName struct {
	Name string `json:"name" validate:"required,correct_text"`
}

type DTOUpdate struct {
	Name    string `json:"name,omitempty" validate:"correct_text"`
	Price   int64  `json:"price,omitempty"`
	Amount  int64  `json:"amount,omitempty"`
	OldName string `json:"oldName" validate:"required,correct_text"`
}
