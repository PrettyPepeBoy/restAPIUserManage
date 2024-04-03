package operations

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"tstUser/internal/http-server/middleware/valid"
	"tstUser/internal/lib/api/response"
	"tstUser/internal/storage/service"
	"tstUser/internal/storage/storages"
	"tstUser/internal/storage/storages/errs"
)

type Response struct {
	response.Response
	user    storages.User
	product storages.Product
}

func BuyProduct(log *slog.Logger, productBuyer service.ProductService, userBuyer service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ps := chi.URLParam(r, "productID")
		productID, err := strconv.Atoi(ps)
		if err != nil {
			log.Error("wrong request productID")
			render.JSON(w, r, response.Error("wrong request productID"))
			return
		}
		ProductID := int64(productID)

		us := chi.URLParam(r, "userID")
		userID, err := strconv.Atoi(us)
		if err != nil {
			log.Error("wrong request userID")
			render.JSON(w, r, response.Error("wrong request userID"))
			return
		}
		UserID := int64(userID)

		user, err := userBuyer.GetUserInfo(UserID)
		if err != nil {
			if errors.Is(err, errs.ErrUserNotFound) {
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
			if errors.Is(err, errs.ErrProductNotFound) {
				log.Error("failed to get product")
				render.JSON(w, r, response.Error("product not found"))
				return
			}
			log.Error("failed to find product")
			render.JSON(w, r, response.Error("failed to find product"))
			return
		}

		err = valid.Buy(user, products)
		if errors.Is(err, errs.ErrProductsEmpty) {
			log.Info("product amount is zero")
			render.JSON(w, r, response.Error("product amount is zero"))
			return
		}
		if errors.Is(err, errs.ErrCashIsNotEnough) {
			log.Info("user's cash is not enough")
			render.JSON(w, r, response.Error("cash is not enough"))
			return
		}

		user.Cash -= products.Price
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

func responseBuyOK(w http.ResponseWriter, r *http.Request, user storages.User, product storages.Product) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		user: storages.User{
			Name:    user.Name,
			Surname: user.Surname,
			Mail:    user.Mail,
			Cash:    user.Cash,
			Date:    user.Date,
			Id:      user.Id,
		},
		product: storages.Product{
			Name:   product.Name,
			Amount: product.Amount,
			Price:  product.Price,
			Id:     product.Id,
		},
	})
}
