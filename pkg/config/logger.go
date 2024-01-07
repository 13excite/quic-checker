package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Defaults logger for the CLI mode
var (
	Base   = zap.NewNop()
	Logger = Base.Sugar()
)

// InitCLILogger loads a global logger for the CLI service
func InitCLILogger() {
	logConfig := zap.NewProductionConfig()
	logConfig.Sampling = nil

	// Log Level
	logConfig.Level.SetLevel(zap.InfoLevel)

	logConfig.Encoding = "console"

	// Set custom encoding format without timestamp
	logConfig.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseColorLevelEncoder,
	}

	logConfig.DisableStacktrace = true
	logConfig.DisableCaller = true

	logConfig.OutputPaths = []string{"stderr"}
	logConfig.ErrorOutputPaths = []string{"stderr"}

	// Build the logger
	globalLogger, _ := logConfig.Build()
	zap.ReplaceGlobals(globalLogger)

	Base = zap.L()
	Logger = Base.Sugar()
}
