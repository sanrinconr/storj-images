// Package controller manage all the entrypoints of the api
package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrResponse test bad responses.
type ErrResponse struct{}

func (e *ErrResponse) Error(*gin.Context) error {
	return NewError(http.StatusInternalServerError, errors.New("test Error"))
}
