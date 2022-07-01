package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping controller to validate health of the app.
type Ping struct{}

// Ping validate the health of the app.
func (p *Ping) Ping(c *gin.Context) error {
	c.Status(http.StatusOK)

	return nil
}
