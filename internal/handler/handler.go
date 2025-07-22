package handler

import (
	"BookStore_API/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
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
	e.GET("/ping", h.serverPing)

	h.registerBookRoutes(e)
}

func (h *Handler) registerBookRoutes(e *echo.Echo) {
	notes := e.Group("/books")
	notes.POST("", h.createBook)
	notes.GET("/:id", h.getByIdBook)
	notes.PUT("/:id", h.updateBook)
	notes.DELETE("/:id", h.deleteBook)
}

func (h *Handler) serverPing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
