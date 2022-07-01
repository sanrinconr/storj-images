package middlewares

import (
	"github.com/gin-gonic/gin"
)

// Bridge make the return error in controller more friendly
// Without bridge return error is c.Error(err), with Bridge is return err.
func Bridge(frontController func(ctx *gin.Context) error) func(*gin.Context) {
	return func(c *gin.Context) {
		if err := frontController(c); err != nil {
			c.Error(err) //nolint:errcheck // if the error is checked, other error is generated
		}
	}
}
