package productDTO

type DTOWithID struct {
	ID     int64  `json:"ID,omitempty"`
	Name   string `json:"name" validate:"required,correct_text"`
	Price  int64  `json:"price" validate:"required"`
	Amount int64  `json:"amount" validate:"required"`
}
