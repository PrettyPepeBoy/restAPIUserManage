package find

import (
	"tstUser/internal/http-server/transport/userDTO"
	"tstUser/internal/lib/api/response"
)

type Request struct {
	ID int `json:"id" validate:"required,id"`
}

type Response struct {
	response.Response
	ID int
	userDTO.UDTO
}

type UserFinder interface {
	GetUserId(ID int)
}
