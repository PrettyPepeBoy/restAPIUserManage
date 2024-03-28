package productDTO

type ProductDTO struct {
	ID     int64  `json:"ID,omitempty"`
	Name   string `json:"name" validate:"required,correct_text"`
	Price  int64  `json:"price" validate:"required"`
	Amount int64  `json:"amount" validate:"required"`
}

type DTOid struct {
	ID int64 `json:"id" validate:"required"`
}

type DTOUpdate struct {
	ID     int64   `json:"id" validate:"required"`
	Name   *string `json:"name,omitempty"`
	Price  *int64  `json:"price,omitempty"`
	Amount *int64  `json:"amount,omitempty"`
}
