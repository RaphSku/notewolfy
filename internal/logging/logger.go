package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupZapLogger(debug bool) *zap.Logger {
	if !debug {
		return zap.NewNop()
	}
	logLevel := zapcore.DebugLevel

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		logLevel,
	)

	loggerOptions := []zap.Option{
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(zap.String("service", "synmake")),
	}

	logger := zap.New(core, loggerOptions...)

	return logger
}
