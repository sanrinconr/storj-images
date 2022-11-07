package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/log"
	"go.uber.org/zap"
)

// LoadLogger set a logger into the context to be latter used.
func LoadLogger(logger *zap.SugaredLogger) func(*gin.Context) {
	return func(ctx *gin.Context) {
		log.AddLoggerToContext(ctx, logger)
		log.SetupUUID(ctx)
		ctx.Next()
	}
}

// Logger is the middleware that generate a log of every request.
func Logger(ctx *gin.Context) {
	ctx.Next()
	// Not used log.Info because the message is custom
	l := log.GetLoggerFromCtx(ctx)
	l.Info(fmt.Sprintf("%s %s %d %s [uuid:%s]",
		ctx.Request.Method,
		ctx.Request.RequestURI,
		ctx.Writer.Status(),
		ctx.Request.Header.Get("X-Real-IP"),
		ctx.Value(log.UUIDKey),
	))
}
