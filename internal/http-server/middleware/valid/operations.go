package valid

import (
	"tstUser/internal/storage/storages"
	"tstUser/internal/storage/storages/errs"
)

func Buy(user storages.User, product storages.Product) error {
	if product.Amount == 0 {
		return errs.ErrProductsEmpty
	}
	if user.Cash-product.Price < 0 {
		return errs.ErrCashIsNotEnough
	}
	return nil
}
