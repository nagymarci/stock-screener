package model

import (
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

var UnknownError = "{\"message\":\"Uknown error\"}"

type HttpError interface {
	Status() int
}

type InternalServerError struct {
	err    string
	status int
}

func NewInternalServerError(msg string) error {
	return &InternalServerError{
		err:    msg,
		status: http.StatusInternalServerError,
	}
}

func (e *InternalServerError) Error() string {
	return e.err
}

func (e *InternalServerError) Status() int {
	return e.status
}

type BadRequestError struct {
	err    string
	status int
}

func NewBadRequestError(msg string) error {
	return &BadRequestError{
		err:    msg,
		status: http.StatusBadRequest,
	}
}

func (e *BadRequestError) Error() string {
	return e.err
}

func (e *BadRequestError) Status() int {
	return e.status
}
