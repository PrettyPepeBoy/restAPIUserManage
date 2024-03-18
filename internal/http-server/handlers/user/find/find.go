package find

import (
	"tstUser/internal/http-server/handlers/user/create"
	"tstUser/internal/lib/api/response"
)

type Request struct {
	ID int `json:"id" validate:"required,id"`
}

type Response struct {
	response.Response
	ID int
	create.UserDTO
}

type UserFinder interface {
	GetUserId(ID int)
}
