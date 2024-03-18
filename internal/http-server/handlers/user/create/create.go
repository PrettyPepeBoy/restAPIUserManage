package create

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"tstUser/internal/http-server/middleware/valid"
	"tstUser/internal/http-server/transport/userDTO"
	resp "tstUser/internal/lib/api/response"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage"
)

type Request struct {
	userDTO.UserDTO
}

type Response struct {
	resp.Response
	Request
}

type UserCreator interface {
	CreateUser(name, surname, date string, cash int) (int64, error)
}

func New(log *slog.Logger, userCreator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server/handlers/user/create/New"
		log = log.With(
			slog.String("op", op),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request body"))
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := valid.CreateValidator().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidateError(validateErr))
			return
		}

		id, err := userCreator.CreateUser(req.Name, req.Surname, req.Date, req.Cash)
		if err != nil {
			if errors.Is(err, storage.ErrUserExist) {
				log.Info("url already exists", slog.String("name", req.Name), slog.String("surname", req.Surname))
				render.JSON(w, r, resp.Error("user already exists"))
				return
			}
			log.Error("failed to add user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add user"))
			return
		}
		log.Info("user added", slog.Int64("id", id))
		responseOK(w, r, req)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, req Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Request: Request{
			userDTO.UserDTO{
				Name:    req.Name,
				Surname: req.Surname,
				Cash:    req.Cash,
				Date:    req.Date,
			},
		},
	})
}
