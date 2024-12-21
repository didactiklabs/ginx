package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitializeLogger(logLevel zapcore.Level) {
	// Create all necessary directories for the log file

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel), // Set the log level
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:      "json",
		EncoderConfig: zap.NewProductionEncoderConfig(),

		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}
	defer Logger.Sync() //nolint:all
}
