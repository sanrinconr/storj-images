package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/src/cmd/dependencies"
	"github.com/sanrinconr/storj-images/src/cmd/middlewares"
)

func setRoutes(router *gin.Engine) {
	resolver := dependencies.New()
	public := router.Group("")
	public.GET("/api/ping", middlewares.Bridge(resolver.Ping()))
	public.POST("/api/image", middlewares.Bridge(resolver.UploadImage()))
	public.GET("/api/image/all", middlewares.Bridge(resolver.GetAllLocations()))

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		panic(err)
	}
}
