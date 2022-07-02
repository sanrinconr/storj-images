package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/controller"
	"go.uber.org/zap"
)

type (
	handler      func(ctx *gin.Context) error
	dependencies struct {
		config Config
		logger *zap.SugaredLogger
	}
)

func resolver() dependencies {
	config := resolveConfig()
	logger := resolveLogger()

	return dependencies{
		config: config,
		logger: logger,
	}
}

// CONTROLLERS

func (d dependencies) Ping() handler {
	ctl := controller.Ping{}

	return ctl.Ping
}

func (d dependencies) Error() handler {
	ctl := controller.DummyError{}

	return ctl.Error
}

// INITIAL CONFIG

func resolveConfig() Config {
	c, err := readConfig()
	if err != nil {
		panic(err)
	}

	return c
}

// LOGGER

func resolveLogger() *zap.SugaredLogger {
	if actualEnvironment() == "prod" {
		return resolveLoggerProd()
	}

	return resolveLoggerDev()
}

func resolveLoggerProd() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	s := logger.Sugar()

	logger.Info("Using logger for production")

	return s
}

func resolveLoggerDev() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	s := logger.Sugar()

	logger.Info("Using logger for development")

	return s
}
