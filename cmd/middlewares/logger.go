package middlewares

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/log"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	return w.ResponseWriter.Write(b)
}

// LoadLogger set a logger into the context to be latter used.
func LoadLogger(logger *zap.SugaredLogger) func(*gin.Context) {
	return func(ctx *gin.Context) {
		log.AddLogger(ctx, logger)
		ctx.Next()
	}
}

// Logger is the middleware that generate a log of every request.
func Logger(ctx *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
	ctx.Writer = blw
	ctx.Next()
	l := log.GetLoggerFromCtx(ctx)
	l.Info(fmt.Sprintf("%s %s %s %d %s",
		ctx.Request.Header.Get("X-Real-IP"),
		ctx.Request.Method,
		ctx.Request.RequestURI,
		ctx.Writer.Status(),
		blw.body.String(),
	))
}
