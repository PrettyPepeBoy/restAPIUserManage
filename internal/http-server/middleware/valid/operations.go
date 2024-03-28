package valid

import (
	"tstUser/internal/http-server/transport/productDTO"
	"tstUser/internal/http-server/transport/userDTO"
	"tstUser/internal/storage"
)

func Buy(user userDTO.UserDTO, product productDTO.ProductDTO) error {
	if product.Amount == 0 {
		return storage.ErrProductsEmpty
	}
	if user.Cash-int(product.Price) < 0 {
		return storage.ErrCashIsNotEnough
	}
	return nil
}
