// Package log manage the creation and generation of custom logs
package log

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	loggerKey  = "loggerKey"
	formatTime = "02-01-2006 03:04:05PM"
)

// New create a new zap logger.
func New(env string) *zap.SugaredLogger {
	logger := resolveLoggerDev()

	if env == "prod" {
		logger = resolveLoggerProd()
	}

	return logger
}

// AddLogger on creation of ctx (used as middleware).
func AddLogger(ctx *gin.Context, l *zap.SugaredLogger) {
	ctx.Set(loggerKey, l)
}

// GetLoggerFromCtx given a context return the saved logger, if not exists, return a default logger.
func GetLoggerFromCtx(ctx context.Context) *zap.SugaredLogger {
	if l, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return l
	}

	// Return default logger dont log anything
	return zap.S()
}

func resolveLoggerProd() *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = sysLogTimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	s := logger.Sugar()

	logger.Info("Using logger for production")

	return s
}

func resolveLoggerDev() *zap.SugaredLogger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = sysLogTimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	s := logger.Sugar()

	logger.Info("Using logger for development")

	return s
}

// sysLogTimeEncoder modify default time log syntax.
func sysLogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(formatTime) + "]")
}
