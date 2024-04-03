package products

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"tstUser/internal/http-server/DTO"
	"tstUser/internal/http-server/middleware/valid"
	"tstUser/internal/lib/api/decode"
	"tstUser/internal/lib/api/response"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage/service"
	"tstUser/internal/storage/storages"
	"tstUser/internal/storage/storages/errs"
)

type Response struct {
	response.Response
	Answer any
}

func CreateProduct(log *slog.Logger, productCreator service.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DTO.ProductDTO
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		product := storages.Product{
			Name:   req.Name,
			Price:  req.Price,
			Amount: req.Amount,
		}
		req.ID, err = productCreator.CreateProducts(product)
		if err != nil {
			if errors.Is(err, errs.ErrProductsExist) {
				log.Info("product already exists", slog.String("name", req.Name))
				render.JSON(w, r, response.Error("product already exists"))
				return
			}
			log.Error("failed to add product", sl.Err(err))
			render.JSON(w, r, response.Error("failed to add product"))
			return
		}
		log.Info("product created", slog.Int64("id", req.ID))
		responseOK(w, r, product)
	}
}

func GetProduct(log *slog.Logger, productGetter service.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DTO.ProductDTOid
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		product, err := productGetter.GetProducts(req.ID)
		if err != nil {
			if errors.Is(err, errs.ErrProductNotFound) {
				log.Info("product doesn't exist", slog.Int64("id", req.ID))
				render.JSON(w, r, response.Error("product doesn't exist"))
				return
			}
			log.Error("failed to get product", sl.Err(err))
			render.JSON(w, r, response.Error("failed to get product"))
			return
		}
		log.Info("got product", slog.Int64("id", product.Id),
			slog.Int("price", product.Price),
			slog.Int("amount", product.Amount),
			slog.String("name", product.Name))
		responseOK(w, r, product)
	}
}

func UpdateProduct(log *slog.Logger, productUpdater service.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DTO.ProductDTOUpdate
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		product, err := productUpdater.GetProducts(req.ID)
		if err != nil {
			if errors.Is(err, errs.ErrProductNotFound) {
				log.Error("there is no such product", errs.ErrProductNotFound)
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
		if err = valid.CreateValidator().Struct(product); err != nil {
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

func DeleteProduct(log *slog.Logger, productDeleter service.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DTO.ProductDTOid
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		err = productDeleter.DeleteProduct(req.ID)
		if err != nil {
			if errors.Is(err, errs.ErrProductNotFound) {
				log.Info("product not found")
				render.JSON(w, r, response.Error("product not found"))
				return
			}
			log.Error("failed to delete product")
			render.JSON(w, r, response.Error("failed to delete product"))
			return
		}
		log.Info("product deleted", slog.Int64("id", req.ID))
		response.OK()
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, req storages.Product) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Answer: DTO.ProductDTO{
			ID:     req.Id,
			Name:   req.Name,
			Price:  req.Price,
			Amount: req.Amount,
		},
	})
}
