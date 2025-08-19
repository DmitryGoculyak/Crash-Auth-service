package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Encoding         string   `mapstructure:"encoding"`
	LogLevel         string   `mapstructure:"level"`
	OutputPaths      []string `mapstructure:"outputPaths"`
	ErrorOutputPaths []string `mapstructure:"errorOutputPaths"`
}

func InitLogger(cfg *Config) (*zap.Logger, error) {
	level, err := zap.ParseAtomicLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	zapCfg := zap.Config{
		Encoding:         cfg.Encoding,
		Level:            zap.NewAtomicLevelAt(level.Level()),
		OutputPaths:      cfg.OutputPaths,
		ErrorOutputPaths: cfg.ErrorOutputPaths,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	logger.Info("Logger initialized", zap.String("level", cfg.LogLevel))
	return logger, nil
}
