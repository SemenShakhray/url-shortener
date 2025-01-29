package response

import (
	"bytes"
	"fmt"

	"github.com/go-playground/validator/v10"
)

const (
	statusOK    = "OK"
	statusError = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func (resp *Response) OK() Response {
	return Response{
		Status: statusOK,
	}
}

func (resp *Response) Err(msg string) Response {
	return Response{
		Status: statusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors, resp Response) Response {
	var buff bytes.Buffer

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			buff.WriteString(fmt.Sprintf("field %s is required field", err.Field()))
		case "url":
			buff.WriteString(fmt.Sprintf("field %s is not valid URL", err.Field()))
		default:
			buff.WriteString(fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return resp.Err(buff.String())
}
