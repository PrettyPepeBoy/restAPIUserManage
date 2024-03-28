package update

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"tstUser/internal/http-server/middleware/valid"
	"tstUser/internal/http-server/transport/userDTO"
	"tstUser/internal/lib/api/decode"
	"tstUser/internal/lib/api/response"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage"
)

type Response struct {
	response.Response
	userDTO.UserDTO
}

type UserUpdater interface {
	UpdateUser(user userDTO.UserDTO) error
	GetUserInfo(ID int64) (userDTO.UserDTO, error)
}

func New(log *slog.Logger, userUpdater UserUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req userDTO.DTOUpdate
		//var user userDTO.UserDTO
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		user, err := userUpdater.GetUserInfo(req.ID)
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user do not exist")
			render.JSON(w, r, response.Error("user do not exist"))
			return
		}
		if err != nil {
			log.Error("failed to find user")
			render.JSON(w, r, response.Error("failed to find user"))
			return
		}
		log.Info("user found", slog.Int64("id", user.ID))
		if req.Name != nil {
			user.Name = *req.Name
		}
		if req.Surname != nil {
			user.Surname = *req.Surname
		}
		if req.Mail != nil {
			user.Mail = *req.Mail
		}
		if req.Cash != nil {
			user.Cash = *req.Cash
		}
		if req.Date != nil {
			user.Cash = *req.Cash
		}
		if err := valid.CreateValidator().Struct(user); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidateError(validateErr))
			return
		}
		err = userUpdater.UpdateUser(user)
		if err != nil {
			log.Error("failed to update user")
			render.JSON(w, r, response.Error("failed to update user"))
			return
		}
		log.Info("user updated", slog.Int64("id", user.ID))

		responseOK(w, r, user)
		//TODO разделить декод чтобы мы могли парсить дефолтные значения
		//TODO апдейт для продукта и юзера
		//TODO покупака товара и история покупок
		//TODO редирект на картинку с товаром
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, user userDTO.UserDTO) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		UserDTO: userDTO.UserDTO{
			ID:      user.ID,
			Name:    user.Name,
			Surname: user.Surname,
			Mail:    user.Mail,
			Cash:    user.Cash,
			Date:    user.Date,
		},
	})
}
