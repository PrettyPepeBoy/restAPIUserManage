package decode

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"tstUser/internal/http-server/middleware/valid"
	resp "tstUser/internal/lib/api/response"
	"tstUser/internal/lib/logger/sl"
)

func Decode(w http.ResponseWriter, r *http.Request, log *slog.Logger, req any) error {
	const op = "lib/api/decode"
	log = log.With(
		slog.String("op", op),
		slog.String("requestID", middleware.GetReqID(r.Context())),
	)
	err := render.DecodeJSON(r.Body, req)
	if err != nil {
		log.Error("Failed to decode request body", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to decode request body"))
		return err
	}

	log.Info("request body decoded", slog.Any("request", req))

	if err := valid.CreateValidator().Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		log.Error("invalid request", sl.Err(err))

		render.JSON(w, r, resp.ValidateError(validateErr))
		return err
	}
	return nil
}
