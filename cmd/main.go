package main

import (
	"BookStore_API/internal/app"
	"BookStore_API/internal/config"
	"BookStore_API/internal/postgres"
	"BookStore_API/internal/zaplog"
	"context"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables.")
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	logger := zaplog.InitLogger(appEnv)
	logger.Info("Logger successfully initialized.")

	logger.Info("Initializing config...")
	cfg := config.InitConfig(logger)
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
