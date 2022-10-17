package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/controller"
	"github.com/sanrinconr/storj-images/cmd/infratructure"
)

type (
	handler      func(ctx *gin.Context) error
	dependencies struct {
		infraStorj infratructure.Storj
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

// INFRA.
func resolveInfraPhotos(c Config) infratructure.Storj {
	t := os.Getenv(c.TokenENV)
	if t == "" {
		panic(fmt.Errorf("variable %s not is defined", c.TokenENV))
	}

	s, err := infratructure.NewStorj(
		infratructure.WithStorjAppAccess(t),
		infratructure.WithStorjBucketName(c.Bucket),
		infratructure.WithStorjProjectName(c.Project),
	)
	if err != nil {
		panic(err)
	}

	return s
}
