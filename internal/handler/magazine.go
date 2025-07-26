package handler

import (
	"BookStore_API/internal/dto"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type CreateMagazineResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}
type GetByIdMagazineResponse struct {
	Magazine dto.MagazineResponse `json:"magazine"`
	Message  string               `json:"message"`
}
type UpdateMagazineResponse struct {
	Message string `json:"message"`
}
type DeleteMagazineResponse struct {
	Message string `json:"message"`
}

func (h *Handler) createMagazine(c echo.Context) error {
	start := time.Now()

	h.logRequestStart(c, "Create magazine request started")

	var req dto.MagazineCreateRequest

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

	magazine := req.ToEntity()

	// create magazine service
	id, err := h.services.Magazine.Create(c.Request().Context(), magazine)
	if err != nil {
		h.logger.Error("failed to create magazine",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusInternalServerError, ErrCreateResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, CreateMagazineResponse{
		Id:      id,
		Message: "magazine created",
	})
}
func (h *Handler) getByIdMagazine(c echo.Context) error {
	start := time.Now()

	h.logRequestStart(c, "Get by id magazine request started")

	// get id param
	id, err := h.parseIdParam(c, start)
	if err != nil {
		return err
	}

	// get by id magazine service
	magazine, err := h.services.Magazine.GetById(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to get by id magazine",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusInternalServerError, ErrGetByIdResponse{
			Message: "internal server error",
		})
	}

	resp := dto.FromEntityMagazine(magazine)

	return c.JSON(http.StatusOK, GetByIdMagazineResponse{
		Magazine: resp,
		Message:  "here is your magazine",
	})
}
func (h *Handler) updateMagazine(c echo.Context) error {
	start := time.Now()

	h.logRequestStart(c, "Update magazine request started")

	// get id param
	id, err := h.parseIdParam(c, start)
	if err != nil {
		return err
	}

	var req dto.MagazineUpdateRequest

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

	// get by id magazine service
	magazine, err := h.services.Magazine.GetById(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to get by id magazine",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusInternalServerError, ErrGetByIdResponse{
			Message: "internal server error",
		})
	}

	req.ApplyToEntity(&magazine)

	// update magazine service
	err = h.services.Magazine.Update(c.Request().Context(), magazine)
	if err != nil {
		h.logger.Error("failed to update magazine",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusInternalServerError, ErrUpdateResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, UpdateMagazineResponse{
		Message: "magazine successfully updated",
	})
}
func (h *Handler) deleteMagazine(c echo.Context) error {
	start := time.Now()

	h.logRequestStart(c, "Delete magazine request started")

	// get id param
	id, err := h.parseIdParam(c, start)
	if err != nil {
		return err
	}

	// delete magazine service
	err = h.services.Magazine.Delete(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to delete by id magazine",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(http.StatusInternalServerError, ErrDeleteByIdResponse{
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, DeleteMagazineResponse{
		Message: "magazine successfully deleted",
	})
}
