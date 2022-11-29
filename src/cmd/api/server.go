// Package api setup the middlewares that are going to be used, the log and configure routes.
package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/src/cmd/config"
	"github.com/sanrinconr/storj-images/src/cmd/middlewares"
	"github.com/sanrinconr/storj-images/src/log"
)

// Start create routes, set port and run app.
func Start() {
	router := createRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		fmt.Printf("%s", err)
	}
}

func createRouter() *gin.Engine {
	logger := log.New(config.ActualEnvironment())
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.LoadLogger(logger))
	router.Use(middlewares.Logger)
	router.Use(middlewares.SetCors())
	router.Use(middlewares.ErrorHandler)
	setRoutes(router)

	return router
}
