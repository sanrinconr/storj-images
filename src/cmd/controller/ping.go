package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewPing create a controller to check the api status.
func NewPing() func(*gin.Context) error {
	return func(c *gin.Context) error {
		c.Status(http.StatusOK)

		return nil
	}
}
