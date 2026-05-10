package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/hatamotoyuki/ninjo/backend/internal/handler/middleware"
	"github.com/hatamotoyuki/ninjo/backend/internal/handler/oapi"
	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

// ListPlans は計画一覧を取得する。
func (h *Handler) ListPlans(ctx echo.Context) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	results, err := h.uc.Plan().List(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	plans := make([]oapi.PlanResponse, len(results))
	for i, r := range results {
		plans[i] = toPlanResponse(r)
	}

	return ctx.JSON(http.StatusOK, oapi.PlanListResponse{Plans: plans})
}

// CreatePlan は計画を作成する。
func (h *Handler) CreatePlan(ctx echo.Context) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	var req oapi.CreatePlanRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid request body"})
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: err.Error()})
	}

	input := usecase.PlanInput{
		Title:      req.Title,
		TotalPages: req.TotalPages,
		TargetDate: req.TargetDate.Time,
	}

	result, err := h.uc.Plan().Create(ctx.Request().Context(), userID, input)
	if err != nil {
		return handlePlanError(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, toPlanResponse(*result))
}

// ReviewPlan は計画をAIにレビューしてもらう。
func (h *Handler) ReviewPlan(ctx echo.Context) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	var req oapi.CreatePlanRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid request body"})
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: err.Error()})
	}

	input := usecase.PlanInput{
		Title:      req.Title,
		TotalPages: req.TotalPages,
		TargetDate: req.TargetDate.Time,
	}

	result, err := h.uc.Plan().Review(ctx.Request().Context(), userID, input)
	if err != nil {
		return handlePlanError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, oapi.PlanReviewResponse{
		DailyPages:    float32(result.DailyPages),
		TotalDays:     result.TotalDays,
		AvailableDays: result.AvailableDays,
		ReviewMessage: result.ReviewMessage,
	})
}

// GetPlan は計画の詳細を取得する。
func (h *Handler) GetPlan(ctx echo.Context, planId openapi_types.UUID) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	result, err := h.uc.Plan().Get(ctx.Request().Context(), userID, planId)
	if err != nil {
		if err == usecase.ErrPlanNotFound {
			return ctx.JSON(http.StatusNotFound, oapi.ErrorResponse{Error: "plan not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	resp := oapi.PlanDetailResponse{
		Id:               result.ID,
		Title:            result.Title,
		TotalPages:       result.TotalPages,
		StartDate:        openapi_types.Date{Time: result.StartDate},
		TargetDate:       openapi_types.Date{Time: result.TargetDate},
		Status:           oapi.PlanDetailResponseStatus(result.Status),
		ProgressRate:     float32(result.ProgressRate),
		DailyPagesNeeded: float32(result.DailyPagesNeeded),
		DaysRemaining:    result.DaysRemaining,
		CreatedAt:        result.CreatedAt,
	}
	if result.AIReview != nil {
		resp.AiReview = result.AIReview
	}

	return ctx.JSON(http.StatusOK, resp)
}

// DeletePlan は計画を削除する。
func (h *Handler) DeletePlan(ctx echo.Context, planId openapi_types.UUID) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	err := h.uc.Plan().Delete(ctx.Request().Context(), userID, planId)
	if err != nil {
		if err == usecase.ErrPlanNotFound {
			return ctx.JSON(http.StatusNotFound, oapi.ErrorResponse{Error: "plan not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	return ctx.NoContent(http.StatusNoContent)
}

func toPlanResponse(r usecase.PlanResult) oapi.PlanResponse {
	resp := oapi.PlanResponse{
		Id:           r.ID,
		Title:        r.Title,
		TotalPages:   r.TotalPages,
		StartDate:    openapi_types.Date{Time: r.StartDate},
		TargetDate:   openapi_types.Date{Time: r.TargetDate},
		Status:       oapi.PlanResponseStatus(r.Status),
		ProgressRate: float32(r.ProgressRate),
		CreatedAt:    r.CreatedAt,
	}
	if r.AIReview != nil {
		resp.AiReview = r.AIReview
	}
	return resp
}

func handlePlanError(ctx echo.Context, err error) error {
	switch err {
	case usecase.ErrInvalidPlan:
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid plan data"})
	case usecase.ErrTargetDatePast:
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "target date must be in the future"})
	case usecase.ErrTooManyPlans:
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "too many active plans (max 10)"})
	default:
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}
}

