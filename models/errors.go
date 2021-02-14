package models

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
)

type responseCode string

type ServiceError struct {
	message    string
	code       responseCode
	httpStatus int
	err        error
}

type serviceErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

const (
	ServiceErrorInvalidRequestBody responseCode = "INVALID_REQUEST_BODY"
	ServiceErrorInternalError      responseCode = "INTERNAL_ERROR"
	ServiceErrorUnauthorised       responseCode = "UNAUTHORISED"
	ServiceErrorUnprocessable           responseCode = "UNPROCESSABLE_ENTITY"
)

func NewServiceError(code responseCode, message string, httpStatus int, err *error) *ServiceError {
	e := errors.New(message)
	if err != nil {
		e = fmt.Errorf("%v: %w", message, *err)
	}
	return &ServiceError{
		message:    message,
		code:       code,
		httpStatus: httpStatus,
		err:        e,
	}
}

func newErrorResponse(code responseCode, message string) *serviceErrorResponse {
	return &serviceErrorResponse{
		Message: message,
		Code:    fmt.Sprintf("%s", code),
	}
}

func (se ServiceError) ToResponse(c echo.Context) error {
	response := newErrorResponse(se.code, se.message)
	return c.JSON(se.httpStatus, response)
}
