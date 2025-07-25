package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

func InitLogger(appEnv string, logFormat string) *zap.Logger {
	var cfg zap.Config

	if appEnv == "production" && logFormat != "console" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.Encoding = "console"
		cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}

	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.StacktraceKey = "stacktrace"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg.InitialFields = map[string]interface{}{
		"service": "bookstore-api",
		"app_env": appEnv,
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatal("Failed to initialize logger: ", err)
	}

	return logger
}
