package app

import (
	"BookStore_API/internal/config"
	"BookStore_API/internal/handler"
	"BookStore_API/internal/repository"
	"BookStore_API/internal/service"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ApplicationRun(cfg *config.Config, logger *zap.Logger, db *pgxpool.Pool) {
	repo := repository.NewRepository(db, logger)
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
