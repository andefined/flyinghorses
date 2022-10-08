package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap.Logger instance
func NewLogger(env, level, path string) *zap.SugaredLogger {
	// Using zap's preset constructors is the simplest way to get a feel for the
	// package, but they don't allow much customization.
	var logger *zap.Logger

	if env == "production" {
		logger = productionLogger(level)
	} else {
		logger = developmentLogger(level)
	}

	defer logger.Sync()
	return logger.Sugar()
}

// set production logger defaults
func productionLogger(level string) *zap.Logger {
	logLevel := zapcore.DebugLevel
	if level == "info" {
		logLevel = zapcore.InfoLevel
	}
	cfg := zap.Config{
		Encoding:         "json",
		Development:      false,
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	l, _ := cfg.Build()
	return l
}

// set development logger defaults
func developmentLogger(level string) *zap.Logger {
	logLevel := zapcore.DebugLevel
	if level == "info" {
		logLevel = zapcore.InfoLevel
	}
	cfg := zap.Config{
		Encoding:         "console",
		Development:      true,
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	l, _ := cfg.Build()
	return l
}
