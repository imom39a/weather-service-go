package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	defaultLogLevel := zapcore.DebugLevel
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(defaultLogLevel)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	output := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(encoder, output, atomicLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.WarnLevel))
	return logger
}
