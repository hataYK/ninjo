package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/hatamotoyuki/ninjo/backend/internal/handler/middleware"
	"github.com/hatamotoyuki/ninjo/backend/internal/handler/oapi"
	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

// GetAvailability は可処分時間設定を取得する。
func (h *Handler) GetAvailability(ctx echo.Context) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	result, err := h.uc.Availability().Get(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	return ctx.JSON(http.StatusOK, toAvailabilityResponse(result))
}

// UpdateAvailability は可処分時間設定を一括更新する。
func (h *Handler) UpdateAvailability(ctx echo.Context) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	var req oapi.UpdateAvailabilityRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid request body"})
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: err.Error()})
	}

	items := make([]usecase.AvailabilityItem, len(req.Availability))
	for i, a := range req.Availability {
		items[i] = usecase.AvailabilityItem{
			DayOfWeek: int8(a.DayOfWeek),
			Hours:     float64(a.Hours),
		}
	}

	result, err := h.uc.Availability().Update(ctx.Request().Context(), userID, items)
	if err != nil {
		if err == usecase.ErrInvalidAvailability {
			return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid availability data"})
		}
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	return ctx.JSON(http.StatusOK, toAvailabilityResponse(result))
}

func toAvailabilityResponse(result *usecase.AvailabilityResult) oapi.AvailabilityResponse {
	items := make([]oapi.AvailabilityItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = oapi.AvailabilityItem{
			DayOfWeek: int(item.DayOfWeek),
			Hours:     float32(item.Hours),
		}
	}
	return oapi.AvailabilityResponse{
		Availability: items,
		WeeklyTotal:  float32(result.WeeklyTotal),
	}
}
