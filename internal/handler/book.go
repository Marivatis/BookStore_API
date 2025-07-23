package handler

import (
	"BookStore_API/internal/dto"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type CreateBookResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}
type GetByIdBookResponse struct {
	Book    dto.BookResponse `json:"book"`
	Message string           `json:"message"`
}

func (h *Handler) createBook(c echo.Context) error {
	start := time.Now()

	h.logRequestStart(c, "Create book request started")

	var req dto.BookCreateRequest

	// request binding
	if err := c.Bind(&req); err != nil {
		h.logger.Error("failed to bind request",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusBadRequest, ErrBindResponse{
			Message: "invalid request body",
		})
	}

	// request validation
	if err := req.Validate(); err != nil {
		h.logger.Error("validation failed",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusBadRequest, ErrValidationResponse{
			Message: err.Error(),
		})
	}

	book := req.ToEntity()

	// create book service
	id, err := h.services.Book.Create(c.Request().Context(), book)
	if err != nil {
		h.logger.Error("failed to create book",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusInternalServerError, ErrCreateResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, CreateBookResponse{
		Id:      id,
		Message: "book created",
	})
}
func (h *Handler) getByIdBook(c echo.Context) error {
	start := time.Now()

	h.logRequestStart(c, "Get by id book request started")

	// get id param
	id, err := h.parseIdParam(c, start)
	if err != nil {
		return err
	}

	// get by id book service
	book, err := h.services.Book.GetById(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to get by id book",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusInternalServerError, ErrGetByIdResponse{
			Message: "internal server error",
		})
	}

	resp := dto.FromEntityBook(book)

	return c.JSON(http.StatusOK, GetByIdBookResponse{
		Book:    resp,
		Message: "here is your book",
	})
}
func (h *Handler) updateBook(c echo.Context) error {
	return nil
}
func (h *Handler) deleteBook(c echo.Context) error {
	return nil
}
