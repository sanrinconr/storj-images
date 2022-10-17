package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/log"
	"github.com/sanrinconr/storj-images/cmd/middlewares"
)

const (
	prod    = "prod"
	develop = "develop"
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
	logger := log.New(actualEnvironment())
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.LoadLogger(logger))
	router.Use(middlewares.Logger)
	router.Use(middlewares.ErrorHandler)
	setRoutes(router)

	return router
}

func actualEnvironment() string {
	env := os.Getenv("ENV")
	if env == "" || env == develop {
		return develop
	}

	return prod
}
