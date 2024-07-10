package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

// ValidationError создает Response с сообщением об ошибке валидации
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("field %s is a required field", err.Field()))
		case "email":
			errMsg = append(errMsg, fmt.Sprintf("field %s is not a valid email address", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, ", "),
	}
}
