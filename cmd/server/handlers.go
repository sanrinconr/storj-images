package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/middlewares"
)

func setRoutes(router *gin.Engine) {
	resolver := resolver(resolveConfig())
	public := router.Group("")
	public.GET("/ping", middlewares.Bridge(resolver.Ping()))
	public.GET("/error", middlewares.Bridge(resolver.Error()))
}

func resolveConfig() Config {
	c, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	return c
}
