package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	statusOK    = "OK"
	statusError = "Error"
)

func OK() Response {
	return Response{
		Status: statusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: statusError,
		Error:  msg,
	}
}

func ValidateError(errs validator.ValidationErrors) Response {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "name":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid name", err.Field()))
		case "surname":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid surname", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid email", err.Field()))
		case "cash":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid cash", err.Field()))
		case "date":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid date", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid", err.Field()))
		}
	}
	return Response{
		Status: statusError,
		Error:  strings.Join(errMsgs, ", ")}
}
