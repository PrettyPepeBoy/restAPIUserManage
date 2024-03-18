package delete

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "tstUser/internal/lib/api/response"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage"
)

type Request struct {
	Id int64 `json:"id"`
}

type Response struct {
	resp.Response
}

type UserDeleter interface {
	DeleteUser(id int64) error
}

func New(log *slog.Logger, userDeleter UserDeleter) http.HandlerFunc {
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

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidateError(validateErr))
			return
		}

		err = userDeleter.DeleteUser(req.Id)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Info("url not exist")
				render.JSON(w, r, resp.Error("user do not exist"))
				return
			}
			log.Error("failed to delete user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete user"))
			return
		}
		log.Info("user deleted", slog.Int64("id", req.Id))

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
