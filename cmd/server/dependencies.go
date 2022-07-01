package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/controller"
)

type (
	handler      func(ctx *gin.Context) error
	dependencies struct{}
)

func resolver() dependencies {
	return dependencies{}
}

// CONTROLLERS

func (d *dependencies) Ping() handler {
	ctl := controller.Ping{}

	return ctl.Ping
}

func (d *dependencies) Error() handler {
	ctl := controller.DummyError{}

	return ctl.Error
}
