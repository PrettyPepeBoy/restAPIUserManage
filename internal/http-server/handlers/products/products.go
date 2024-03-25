package products

import (
	"errors"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tstUser/internal/http-server/transport/productDTO"
	"tstUser/internal/lib/api/decode"
	"tstUser/internal/lib/api/response"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage"
)

type Request struct {
	productDTO.DTOWithID
}

type Response struct {
	response.Response
	RequestAnswer Request
}

type ProductCreator interface {
	CreateProducts(name string, price, amount int64) (int64, error)
}

type ProductUpdater interface {
	//TODO implement method
}

func CreateProduct(log *slog.Logger, productCreator ProductCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		req.ID, err = productCreator.CreateProducts(req.Name, req.Price, req.Amount)
		if err != nil {
			if errors.Is(err, storage.ErrProductsExist) {
				log.Info("product already exists", slog.String("name", req.Name))
				render.JSON(w, r, response.Error("product already exists"))
				return
			}
			log.Error("failed to add product", sl.Err(err))
			render.JSON(w, r, response.Error("failed to add product"))
			return
		}
		log.Info("product created", slog.Int64("id", req.ID))
		responseCreateProductOK(w, r, req)
	}
}

//func UpdateProduct(log *slog.Logger, pr)

func responseCreateProductOK(w http.ResponseWriter, r *http.Request, req Request) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		RequestAnswer: Request{
			productDTO.DTOWithID{
				ID:     req.ID,
				Name:   req.Name,
				Price:  req.Price,
				Amount: req.Amount,
			},
		},
	})
}
