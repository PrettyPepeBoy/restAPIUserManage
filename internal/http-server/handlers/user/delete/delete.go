package delete

import (
	"errors"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tstUser/internal/lib/api/decode"
	resp "tstUser/internal/lib/api/response"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage"
)

type Request struct {
	ID int64 `json:"id"`
}

type Response struct {
	resp.Response
	Request
}

type UserDeleter interface {
	DeleteUser(id int64) error
}

func New(log *slog.Logger, userDeleter UserDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}

		err = userDeleter.DeleteUser(req.ID)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Info("user not exist")
				render.JSON(w, r, resp.Error("user do not exist"))
				return
			}
			log.Error("failed to delete user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete user"))
			return
		}
		log.Info("user deleted", slog.Int64("id", req.ID))

		responseOK(w, r, req.ID)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, req int64) {
	render.JSON(w, r, Response{
		Request:  Request{ID: req},
		Response: resp.OK(),
	})
}
