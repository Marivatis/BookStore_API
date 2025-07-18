package handler

import (
	"BookStore_API/internal/entity"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) createBook(c echo.Context) error {
	var book entity.Book

	// читаємо тіло запиту в структуру
	if err := c.Bind(&book); err != nil {
		h.logger.Error("failed to bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request payload",
		})
	}

	// викликаємо сервіс
	id, err := h.services.Book.Create(c.Request().Context(), book)
	if err != nil {
		h.logger.Error("failed to create book", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create book",
		})
	}

	// повертаємо успіх
	return c.JSON(http.StatusCreated, map[string]any{
		"id":     id,
		"status": "book created",
	})
}
