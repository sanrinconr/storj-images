// Package log manage the creation and generation of custom logs
package log

import (
	"context"
	"crypto/rand"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// UUIDKey Is used to find the uuid of a request into the context.
	UUIDKey    = "uuidKey"
	loggerKey  = "loggerKey"
	formatTime = "02-01-2006 03:04:05PM"
)

var absolutePath = getAbsolutePath()

// New create a new zap logger.
func New(env string) *zap.SugaredLogger {
	logger := resolveLoggerDev()

	if env == "prod" {
		logger = resolveLoggerProd()
	}

	logger.Info(fmt.Sprintf("Using '%s' environment", env))

	return logger
}

// AddLoggerToContext on creation of ctx (used as middleware).
func AddLoggerToContext(ctx *gin.Context, l *zap.SugaredLogger) {
	ctx.Set(loggerKey, l)
}

// SetupUUID generate and load a request id into the context.
func SetupUUID(ctx *gin.Context) {
	ctx.Set(UUIDKey, generateID())
}

// Debug generate a log of level debug.
func Debug(ctx context.Context, msg string) {
	uuid := ctx.Value(UUIDKey)
	GetLoggerFromCtx(ctx).Debug(formatMsg(uuid, msg))
}

// Info generate a log of level info.
func Info(ctx context.Context, msg string) {
	uuid := ctx.Value(UUIDKey)
	GetLoggerFromCtx(ctx).Info(formatMsg(uuid, msg))
}

// Error generate a log of level error.
func Error(ctx context.Context, err error) {
	uuid := ctx.Value(UUIDKey)
	GetLoggerFromCtx(ctx).Error(formatMsg(uuid, err.Error()))
}

// GetLoggerFromCtx given a context return the saved logger, if not exists, return a default logger.
func GetLoggerFromCtx(ctx context.Context) *zap.SugaredLogger {
	if l, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return l
	}

	// Return default logger dont log anything
	return zap.S()
}

func formatMsg(uuid any, msg string) string {
	return fmt.Sprintf("[msg:%s] [uuid:%s] [%s] ", msg, uuid, lastCaller())
}

func lastCaller() string {
	const stackCallsBehind = 3
	_, file, line, _ := runtime.Caller(stackCallsBehind)

	return fmt.Sprintf("%s:%d", relative(file), line)
}

func relative(path string) string {
	return strings.TrimPrefix(filepath.ToSlash(path), absolutePath)
}

func getAbsolutePath() string {
	_, fileName, _, _ := runtime.Caller(0)

	return filepath.ToSlash(filepath.Dir(filepath.Dir(fileName))) + "/"
}

func resolveLoggerProd() *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = sysLogTimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}

func resolveLoggerDev() *zap.SugaredLogger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = sysLogTimeEncoder
	cfg.DisableCaller = true

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}

// sysLogTimeEncoder modify default time log syntax.
func sysLogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(formatTime) + "]")
}

func generateID() string {
	const len = 8 // https://en.wikipedia.org/wiki/Universally_unique_identifier

	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", b)
}
