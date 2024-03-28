package products

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"tstUser/internal/http-server/middleware/valid"
	"tstUser/internal/http-server/transport/productDTO"
	"tstUser/internal/lib/api/decode"
	"tstUser/internal/lib/api/response"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage"
)

type Response struct {
	response.Response
	Answer any
}

type ProductCreator interface {
	CreateProducts(name string, price, amount int64) (int64, error)
}

type ProductGetter interface {
	GetProducts(ID int64) (productDTO.ProductDTO, error)
}

type ProductUpdater interface {
	GetProducts(ID int64) (productDTO.ProductDTO, error)
	UpdateProducts(up productDTO.ProductDTO) error
}

func CreateProduct(log *slog.Logger, productCreator ProductCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req productDTO.ProductDTO
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
		responseOK(w, r, req)
	}
}

func GetProduct(log *slog.Logger, productGetter ProductGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req productDTO.DTOid
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		product, err := productGetter.GetProducts(req.ID)
		if err != nil {
			if errors.Is(err, storage.ErrProductNotFound) {
				log.Info("product doesn't exist", slog.Int64("id", req.ID))
				render.JSON(w, r, response.Error("product doesn't exist"))
				return
			}
			log.Error("failed to get product", sl.Err(err))
			render.JSON(w, r, response.Error("failed to get product"))
			return
		}
		if product.Amount == 0 {
			log.Error("product is empty", storage.ErrProductsEmpty)
			render.JSON(w, r, response.Error("product is empty"))
			return
		}
		log.Info("got product", slog.Int64("id", product.ID),
			slog.Int64("price", product.Price),
			slog.Int64("amount", product.Amount),
			slog.String("name", product.Name))
		responseOK(w, r, product)
	}
}

func UpdateProduct(log *slog.Logger, productUpdater ProductUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req productDTO.DTOUpdate
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		product, err := productUpdater.GetProducts(req.ID)
		if err != nil {
			if errors.Is(err, storage.ErrProductNotFound) {
				log.Error("there is no such product", storage.ErrProductNotFound)
				render.JSON(w, r, response.Error("product not found"))
				return
			}
			log.Error("failed to get product", err)
			render.JSON(w, r, response.Error("failed to get product"))
			return
		}
		if req.Name != nil {
			product.Name = *req.Name
		}
		if req.Amount != nil {
			product.Amount = *req.Amount
		}
		if req.Price != nil {
			product.Price = *req.Price
		}
		if err := valid.CreateValidator().Struct(product); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidateError(validateErr))
			return
		}
		err = productUpdater.UpdateProducts(product)
		if err != nil {
			log.Error("failed to update product", sl.Err(err))
			render.JSON(w, r, response.Error("failed to update product"))
			return
		}
		log.Info("product updated")
		response.OK()
	}
}

//TODO ОПЕРАЦИИ

func responseOK(w http.ResponseWriter, r *http.Request, req productDTO.ProductDTO) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Answer: productDTO.ProductDTO{
			ID:     req.ID,
			Name:   req.Name,
			Price:  req.Price,
			Amount: req.Amount,
		},
	})
}
