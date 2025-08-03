package handler

import (
	"BookStore_API/internal/dto"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type CreateOrderResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}
type GetByIdOrderResponse struct {
	Order   dto.OrderResponse `json:"order"`
	Message string            `json:"message"`
}
type UpdateOrderResponse struct {
	Message string `json:"message"`
}
type DeleteOrderResponse struct {
	Message string `json:"message"`
}

func (h *Handler) createOrder(c echo.Context) error {
	start := time.Now()

	h.logRequestStart(c, "Create book request started")

	var req dto.OrderCreateRequest

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

	order := req.ToEntity()

	// create book service
	id, err := h.services.Order.Create(c.Request().Context(), order)
	if err != nil {
		h.logger.Error("failed to create order",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusInternalServerError, ErrCreateResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, CreateOrderResponse{
		Id:      id,
		Message: "order created",
	})
}
func (h *Handler) getByIdOrder(c echo.Context) error {
	start := time.Now()

	h.logRequestStart(c, "Get by id order request started")

	// get id param
	id, err := h.parseIdParam(c, start)
	if err != nil {
		return err
	}

	// get by id order service
	order, err := h.services.Order.GetById(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to get by id order",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusInternalServerError, ErrGetByIdResponse{
			Message: "internal server error",
		})
	}

	resp := dto.FromEntityOrder(order)

	return c.JSON(http.StatusOK, GetByIdOrderResponse{
		Order:   resp,
		Message: "here is your order",
	})
}
func (h *Handler) updateOrder(c echo.Context) error {
	start := time.Now()

	h.logRequestStart(c, "Update order request started")

	// get id param
	id, err := h.parseIdParam(c, start)
	if err != nil {
		return err
	}

	var req dto.OrderUpdateRequest

	// request binding
	if err = c.Bind(&req); err != nil {
		h.logger.Error("failed to bind request",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusBadRequest, ErrBindResponse{
			Message: "invalid request body",
		})
	}

	// request validation
	if err = req.Validate(); err != nil {
		h.logger.Error("validation failed",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusBadRequest, ErrValidationResponse{
			Message: err.Error(),
		})
	}

	// get by id order service
	order, err := h.services.Order.GetById(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to get by id order",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusInternalServerError, ErrGetByIdResponse{
			Message: "internal server error",
		})
	}

	req.ApplyToEntity(&order)

	// update order service
	err = h.services.Order.Update(c.Request().Context(), order)
	if err != nil {
		h.logger.Error("failed to update order",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusInternalServerError, ErrUpdateResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, UpdateOrderResponse{
		Message: "order successfully updated",
	})
}
func (h *Handler) deleteOrder(c echo.Context) error {
	start := time.Now()

	h.logRequestStart(c, "Delete order request started")

	// get id param
	id, err := h.parseIdParam(c, start)
	if err != nil {
		return err
	}

	// delete order service
	err = h.services.Order.Delete(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to delete by id book",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusInternalServerError, ErrDeleteByIdResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, DeleteOrderResponse{
		Message: "order successfully deleted",
	})
}
