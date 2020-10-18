package model

type ErrorResponse struct {
	Message string `json:"message"`
}

var UnknownError = "{\"message\":\"Uknown error\"}"
