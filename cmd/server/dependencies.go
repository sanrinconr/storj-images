package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/controller"
	"github.com/sanrinconr/storj-images/cmd/infrastructure"
)

type (
	handler      func(ctx *gin.Context) error
	dependencies struct {
		infraStorj infrastructure.Storj
	}
)

func resolver(c Config) dependencies {
	d := dependencies{
		infraStorj: resolveInfraPhotos(c),
	}

	return d
}

// CONTROLLERS

func (d dependencies) Ping() handler {
	ctl := controller.Ping{}

	return ctl.Ping
}

func (d dependencies) Error() handler {
	ctl := controller.ErrResponse{}

	return ctl.Error
}

// STORJ

func resolveInfraPhotos(c Config) infrastructure.Storj {
	t := os.Getenv(c.Database.Photos.TokenENV)
	if t == "" {
		panic(fmt.Errorf("variable %s not exists", c.Database.Photos.TokenENV))
	}

	s, err := infrastructure.NewStorj(
		infrastructure.WithStorjAppAccess(t),
		infrastructure.WithStorjBucketName(c.Database.Photos.Bucket),
		infrastructure.WithStorjProjectName(c.Database.Photos.Project),
	)
	if err != nil {
		panic(err)
	}

	return s
}
