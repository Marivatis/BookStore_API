package app

import (
	"BookStore_API/internal/config"
	"BookStore_API/internal/handler"
	"BookStore_API/internal/repository"
	"BookStore_API/internal/service"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ApplicationRun(cfg *config.Config, logger *zap.Logger, db *sql.DB) {
	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services, logger)

	runServer(handlers, cfg.Port, logger)
}

func runServer(h *handler.Handler, port int, logger *zap.Logger) {
	e := echo.New()

	h.RegisterRoutes(e)

	addr := fmt.Sprintf(":%d", port)
	logger.Info("Starting server...", zap.String("address", addr))

	if err := e.Start(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
