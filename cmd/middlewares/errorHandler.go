package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/controller"
)

// ErrorHandler find if in the context exists errors, if yes return a formated JSON.
func ErrorHandler(ctx *gin.Context) {
	ctx.Next()

	if len(ctx.Errors) > 0 {
		err := ctx.Errors.Last()

		var typed controller.Error

		if errors.As(err, &typed) {
			ctx.JSON(typed.Code, typed)
		} else {
			entity := controller.Error{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Cause:   err.Error(),
			}
			ctx.JSON(entity.Code, entity)
		}
	}
}
