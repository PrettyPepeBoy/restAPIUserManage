package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tstUser/internal/lib/api/response"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage"
)

type UserFinder interface {
	GetUser(userMail string) error
}

func New(log *slog.Logger, userFinder UserFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/redirect/New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		mail := chi.URLParam(r, "mail")
		if mail == "" {
			log.Info("mail is empty")
			render.JSON(w, r, response.Error("not found"))
			return
		}
		err := userFinder.GetUser(mail)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Info("mail not found", "mail", mail)
				render.JSON(w, r, response.Error("not found"))
				return
			}
			log.Error("failed to get mail", sl.Err(err))
			render.JSON(w, r, response.Error("not found"))
			return
		}
		log.Info("got mail", slog.String("mail", mail))
		http.Redirect(w, r, mail, http.StatusFound)
	}
}
