package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/middlewares"
)

func setRoutes(router *gin.Engine) {
	resolver := resolver()
	public := router.Group("")
	public.GET("/ping", middlewares.Bridge(resolver.Ping()))
	public.GET("/error", middlewares.Bridge(resolver.Error()))
}
