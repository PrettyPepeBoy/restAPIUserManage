package check

import (
	"errors"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tstUser/internal/lib/api/decode"
	"tstUser/internal/lib/api/response"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage"
)

type Request struct {
	ID int64 `json:"id"`
}

type Response struct {
	response.Response
	ID   int64
	Mail string
}

type UserChecker interface {
	CheckUserID(ID int64) (string, error)
}

func New(log *slog.Logger, userFinder UserChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		userMail, err := userFinder.CheckUserID(req.ID)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Info("can not find user with such id", slog.Int64("id", req.ID))
				render.JSON(w, r, response.Error("user not found"))
				return
			}
			log.Error("failed to find user", sl.Err(err))
			render.JSON(w, r, response.Error("failed to find user"))
			return
		}
		log.Info("user founded", slog.Int64("id", req.ID))
		responseOK(w, r, req, userMail)
		http.Redirect(w, r, userMail, http.StatusFound)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, req Request, userMail string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		ID:       req.ID,
		Mail:     userMail,
	})
}
