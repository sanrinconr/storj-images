package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/middlewares"
	"github.com/sanrinconr/storj-images/cmd/server/dependencies"
)

func setRoutes(router *gin.Engine) {
	resolver := dependencies.New(actualEnvironment())
	public := router.Group("")
	public.GET("/ping", middlewares.Bridge(resolver.Ping()))
	public.GET("/error", middlewares.Bridge(resolver.Error()))
	public.POST("/image", middlewares.Bridge(resolver.AddImage()))

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		panic(err)
	}
}
