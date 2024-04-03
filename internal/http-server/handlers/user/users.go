package user

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"tstUser/internal/http-server/DTO"
	"tstUser/internal/http-server/middleware/valid"
	"tstUser/internal/lib/api/decode"
	resp "tstUser/internal/lib/api/response"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage/service"
	"tstUser/internal/storage/storages"
	"tstUser/internal/storage/storages/errs"
)

type Response struct {
	resp.Response
	Answer any
}

func CreateUser(log *slog.Logger, userCreator service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DTO.UserDTO
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		user := storages.User{
			Name:    req.Name,
			Surname: req.Surname,
			Mail:    req.Mail,
			Cash:    req.Cash,
			Date:    req.Date,
		}
		id, err := userCreator.CreateUser(user)
		if err != nil {
			if errors.Is(err, errs.ErrUserExist) {
				log.Info("user already exists", slog.String("mail", req.Mail))
				render.JSON(w, r, resp.Error("user already exists"))
				return
			}
			log.Error("failed to add user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add user"))
			return
		}
		log.Info("user added", slog.Int64("id", id))
		user.Id = id
		responseOK(w, r, user)
	}
}

func DeleteUser(log *slog.Logger, userDeleter service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DTO.UserDTO
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		err = userDeleter.DeleteUser(req.ID)
		if err != nil {
			if errors.Is(err, errs.ErrUserNotFound) {
				log.Info("user not exist")
				render.JSON(w, r, resp.Error("user do not exist"))
				return
			}
			log.Error("failed to delete user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete user"))
			return
		}
		log.Info("user deleted", slog.Int64("id", req.ID))
		resp.OK()
	}
}

func UpdateUser(log *slog.Logger, userUpdater service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DTO.UserDTOUpdate
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		user, err := userUpdater.GetUserInfo(req.ID)
		if errors.Is(err, errs.ErrUserNotFound) {
			log.Error("user do not exist")
			render.JSON(w, r, resp.Error("user do not exist"))
			return
		}
		if err != nil {
			log.Error("failed to find user")
			render.JSON(w, r, resp.Error("failed to find user"))
			return
		}
		log.Info("user found", slog.Int64("id", user.Id))
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
		if err = valid.CreateValidator().Struct(user); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidateError(validateErr))
			return
		}
		err = userUpdater.UpdateUser(user)
		if err != nil {
			log.Error("failed to update user")
			render.JSON(w, r, resp.Error("failed to update user"))
			return
		}
		log.Info("user updated", slog.Int64("id", user.Id))
		responseOK(w, r, user)
	}
}

func FindUser(log *slog.Logger, userFinder service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DTO.UserDTOid
		err := decode.Decode(w, r, log, &req)
		if err != nil {
			return
		}
		user, err := userFinder.GetUserInfo(req.Id)
		if err != nil {
			if errors.Is(err, errs.ErrUserNotFound) {
				log.Info("user is not exist")
				render.JSON(w, r, resp.Error("user is not exist"))
				return
			}
			log.Error("failed to find user", err)
			render.JSON(w, r, resp.Error("failed to find user"))
			return
		}
		log.Info("user found", slog.Int64("id", user.Id))
		responseOK(w, r, user)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, user storages.User) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Answer: storages.User{
			Id:      user.Id,
			Name:    user.Name,
			Surname: user.Surname,
			Mail:    user.Mail,
			Cash:    user.Cash,
			Date:    user.Date,
		},
	})
}
