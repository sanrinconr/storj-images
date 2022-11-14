package controller

import (
	"net/http"
)

// Error returned when a error is generated, can be used by all controllers.
type Error struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
	Cause   string `json:"Cause"`
}

// NewError controller that return a error.
func NewError(status int, err error) Error {
	return Error{
		Code:    status,
		Message: http.StatusText(status),
		Cause:   err.Error(),
	}
}

func (e Error) Error() string {
	return e.Cause
}
