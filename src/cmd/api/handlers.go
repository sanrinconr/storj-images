package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/src/cmd/dependencies"
	"github.com/sanrinconr/storj-images/src/cmd/middlewares"
)

func setRoutes(router *gin.Engine) {
	resolver := dependencies.New()
	public := router.Group("")
	public.GET("/ping", middlewares.Bridge(resolver.Ping()))
	public.POST("/image", middlewares.Bridge(resolver.UploadImage()))

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		panic(err)
	}
}
