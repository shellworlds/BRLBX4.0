package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a production or development logger.
func New(dev bool) (*zap.Logger, error) {
	if dev {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return cfg.Build()
	}
	return zap.NewProduction()
}
