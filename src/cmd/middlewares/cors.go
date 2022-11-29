package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/src/cmd/config"
)

// SetCors configure the headers to allow CORS from the front-end.
func SetCors() gin.HandlerFunc {
	allowedOrigins := []string{"*"}
	if env := config.ActualEnvironment(); env == config.Prod {
		allowedOrigins = []string{"https://photos.santiagorincon.tk"}
	}

	return cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    allowedOrigins,
		AllowMethods:    []string{"GET"},
		AllowHeaders:    []string{"Origin"},
		ExposeHeaders:   []string{"Content-Length"},
		AllowWildcard:   false,
	})
}
