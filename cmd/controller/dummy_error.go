package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DummyError test bad responses.
type DummyError struct{}

func (d DummyError) Error(*gin.Context) error {
	return NewError(http.StatusInternalServerError, errors.New("dummy error"))
}
