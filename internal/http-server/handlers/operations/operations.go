package operations

import (
	"errors"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tstUser/internal/http-server/middleware/valid"
	"tstUser/internal/http-server/transport/productDTO"
	"tstUser/internal/http-server/transport/userDTO"
	"tstUser/internal/lib/api/response"
	"tstUser/internal/storage"
)

type Response struct {
	response.Response
	userDTO.UserDTO
	productDTO.ProductDTO
}

type ProductBuyer interface {
	GetProducts(ID int64) (productDTO.ProductDTO, error)
	UpdateProducts(up productDTO.ProductDTO) error
}

type UserBuyer interface {
	GetUserInfo(ID int64) (userDTO.UserDTO, error)
	UpdateUser(user userDTO.UserDTO) error
}

func BuyProduct(log *slog.Logger, ProductID, UserID int64, productBuyer ProductBuyer, userBuyer UserBuyer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userBuyer.GetUserInfo(UserID)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Error("user is not exist")
				render.JSON(w, r, response.Error("user is not exist"))
				return
			}
			log.Error("failed to find user")
			render.JSON(w, r, response.Error("failed to find user"))
			return
		}
		products, err := productBuyer.GetProducts(ProductID)
		if err != nil {
			if errors.Is(err, storage.ErrProductNotFound) {
				log.Error("failed to get product")
				render.JSON(w, r, response.Error("product not found"))
				return
			}
			log.Error("failed to find product")
			render.JSON(w, r, response.Error("failed to find product"))
			return
		}

		err = valid.Buy(user, products)
		if errors.Is(err, storage.ErrProductsEmpty) {
			log.Info("product amount is zero")
			render.JSON(w, r, response.Error("product amount is zero"))
			return
		}
		if errors.Is(err, storage.ErrCashIsNotEnough) {
			log.Info("user's cash is not enough")
			render.JSON(w, r, response.Error("cash is not enough"))
			return
		}

		user.Cash -= int(products.Price)
		err = userBuyer.UpdateUser(user)
		if err != nil {
			log.Error("failed to update User")
			return
		}
		products.Amount -= 1

		err = productBuyer.UpdateProducts(products)
		if err != nil {
			log.Error("failed to update product")
			return
		}
		log.Info("success")
		responseBuyOK(w, r, user, products)
		//TODO нужно как то сделать, чтобы данные юзера и продукта обновлялись одновременно
	}
}

func responseBuyOK(w http.ResponseWriter, r *http.Request, user userDTO.UserDTO, products productDTO.ProductDTO) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		UserDTO: userDTO.UserDTO{
			Name:    user.Name,
			Surname: user.Surname,
			Mail:    user.Mail,
			Cash:    user.Cash,
			Date:    user.Date,
			ID:      user.ID,
		},
		ProductDTO: productDTO.ProductDTO{
			Name:   products.Name,
			Amount: products.Amount,
			Price:  products.Price,
			ID:     products.ID,
		},
	})
}
