package create

import (
	"errors"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tstUser/internal/http-server/transport/userDTO"
	"tstUser/internal/lib/api/decode"
	resp "tstUser/internal/lib/api/response"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage"
)

type Request struct {
	userDTO.UDTO
}

type Response struct {
	resp.Response
	Request
	ID int64
}

type UserCreator interface {
	CreateUser(name, surname, mail, date string, cash int) (int64, error)
}

func New(log *slog.Logger, userCreator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		id, err := userCreator.CreateUser(req.Name, req.Surname, req.Mail, req.Date, req.Cash)
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
		responseOK(w, r, req, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, req Request, id int64) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		ID:       id,
		Request: Request{
			userDTO.UDTO{
				Name:    req.Name,
				Surname: req.Surname,
				Mail:    req.Mail,
				Cash:    req.Cash,
				Date:    req.Date,
			},
		},
	})
}
