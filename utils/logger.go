package utils

import (
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// Logger instance
var Logger *zap.Logger

// SetupLogger initialize a zap logger instance.
func SetupLogger(l *lumberjack.Logger) *zap.Logger {

	if !l.Compress {
		l.Compress = true
	}

	w := zapcore.AddSync(l)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)

	Logger = zap.New(core)

	return Logger
}

// Info logs request info with pre-defined fields.
func Info(r *http.Request, s string) {
	fields := []zapcore.Field{
		zap.String("remote", r.RemoteAddr),
		zap.String("request", r.RequestURI),
		zap.String("method", r.Method),
	}

	Logger.Info(s, fields...)
}

// Error logs error with pre-defined fields.
func Error(r *http.Request, err error) {
	fields := []zapcore.Field{
		zap.String("remote", r.RemoteAddr),
		zap.String("request", r.RequestURI),
		zap.String("method", r.Method),
	}

	Logger.Error(err.Error(), fields...)
}
