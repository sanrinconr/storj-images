package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/controller"
	"github.com/sanrinconr/storj-images/cmd/infratructure"
	"go.uber.org/zap"
)

type (
	handler      func(ctx *gin.Context) error
	dependencies struct {
		config     Config
		logger     *zap.SugaredLogger
		infraStorj infratructure.Storj
	}
)

func resolver() dependencies {
	config := resolveConfig()
	logger := resolveLogger()
	infraPhotos := resolveInfraPhotos(config, logger)

	return dependencies{
		config:     config,
		logger:     logger,
		infraStorj: infraPhotos,
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

// STORJ

func resolveInfraPhotos(c Config, l *zap.SugaredLogger) infratructure.Storj {
	t := os.Getenv(c.TokenENV)
	if t == "" {
		panic(fmt.Errorf("variable %s not exists", c.TokenENV))
	}

	s, err := infratructure.NewStorj(
		infratructure.WithStorjAppAccess(t),
		infratructure.WithStorjBucketName(c.Bucket),
		infratructure.WithStorjProjectName(c.Project),
		infratructure.WithStorjLogger(l),
	)
	if err != nil {
		panic(err)
	}

	return s
}

// INITIAL CONFIG.
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
