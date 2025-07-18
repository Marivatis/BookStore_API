package handler

import (
	"BookStore_API/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	services *service.Service
	logger   *zap.Logger
}

func NewHandler(s *service.Service, logger *zap.Logger) *Handler {
	return &Handler{
		services: s,
		logger:   logger,
	}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.POST("/books", h.createBook)
}
