package main

import (
	"BookStore_API/internal/app"
	"BookStore_API/internal/config"
	"BookStore_API/internal/postgres"
	"BookStore_API/internal/zaplog"
	"context"
	"go.uber.org/zap"
	"os"
	"time"
)

func main() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}
	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "" {
		logFormat = "json"
	}

	logger := zaplog.InitLogger(appEnv, logFormat)
	logger.Info("Logger successfully initialized.")

	logger.Info("Initializing config...")
	cfg := config.InitConfig(logger, appEnv)
	logger.Info("Config successfully initialized:", zap.Any("cfg", cfg))

	logger.Info("Initializing DB connection...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := postgres.NewPostgresDB(ctx, &cfg.DBCfg)
	if err != nil {
		logger.Fatal("Failed to initialize DB connection.", zap.Error(err))
	}
	defer func() {
		logger.Info("Closing DB connection...")
		db.Close()
	}()
	logger.Info("DB connection successfully initialized.")

	logger.Info("Running application...")
	app.ApplicationRun(cfg, logger, db)
	logger.Info("Application closed.")
}
