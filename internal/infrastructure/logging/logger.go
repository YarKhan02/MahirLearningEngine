package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Env 		string
	ServiceName string
	Version 	string
	LogFilePath	string
}

func New(cfg Config) (*zap.Logger, error) {
	var encoderCfg zapcore.EncoderConfig
	var level zapcore.Level

	if cfg.Env == "development" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		level = zapcore.DebugLevel
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.LevelKey = "level"
		encoderCfg.MessageKey = "message"
		level = zapcore.InfoLevel
	}

	var writer zapcore.WriteSyncer
	var encoder zapcore.Encoder

	if cfg.Env == "development" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
		logFile, err := os.OpenFile(cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.LogFilePath,
			MaxSize:    100,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		})
		_ = logFile
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg) // JSON in prod — Alloy/Loki parse this cleanly
		writer = zapcore.AddSync(os.Stdout) // Render captures container stdout natively
	}

	core := zapcore.NewCore(encoder, writer, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Fields attached to every single log line
	logger = logger.With(
		zap.String("service", 	cfg.ServiceName),
		zap.String("version", 	cfg.Version),
		zap.String("env", 		cfg.Env),
	)

	return logger, nil
}