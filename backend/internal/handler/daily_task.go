package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/hatamotoyuki/ninjo/backend/internal/handler/middleware"
	"github.com/hatamotoyuki/ninjo/backend/internal/handler/oapi"
	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

// ListDailyTasks は指定日のデイリータスク一覧を取得する。
func (h *Handler) ListDailyTasks(ctx echo.Context, params oapi.ListDailyTasksParams) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	date, err := time.Parse("2006-01-02", params.Date.Format("2006-01-02"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid date format"})
	}

	result, err := h.uc.DailyTask().List(ctx.Request().Context(), userID, date)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	return ctx.JSON(http.StatusOK, toDailyTaskListResponse(result))
}

// GenerateDailyTasks は指定日のデイリータスクを自動生成する。
func (h *Handler) GenerateDailyTasks(ctx echo.Context) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	var req oapi.GenerateDailyTasksRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid request body"})
	}

	date, err := time.Parse("2006-01-02", req.Date.Format("2006-01-02"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid date format"})
	}

	result, err := h.uc.DailyTask().Generate(ctx.Request().Context(), userID, date)
	if err != nil {
		return handleDailyTaskError(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, toDailyTaskListResponse(result))
}

// UpdateDailyTask はタスクのページ範囲を更新する。
func (h *Handler) UpdateDailyTask(ctx echo.Context, taskId openapi_types.UUID) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	var req oapi.UpdateDailyTaskRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid request body"})
	}

	result, err := h.uc.DailyTask().Update(ctx.Request().Context(), userID, taskId, req.StartPage, req.EndPage)
	if err != nil {
		return handleDailyTaskError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, toDailyTaskResponse(*result))
}

// CompleteDailyTask はタスクを完了にする。
func (h *Handler) CompleteDailyTask(ctx echo.Context, taskId openapi_types.UUID) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	var req oapi.CompleteDailyTaskRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid request body"})
	}

	result, err := h.uc.DailyTask().Complete(ctx.Request().Context(), userID, taskId, req.ActualEndPage, req.Memo)
	if err != nil {
		return handleDailyTaskError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, toDailyTaskResponse(*result))
}

func toDailyTaskResponse(r usecase.DailyTaskResult) oapi.DailyTaskResponse {
	resp := oapi.DailyTaskResponse{
		Id:          r.ID,
		PlanId:      r.PlanID,
		PlanTitle:   r.PlanTitle,
		Date:        openapi_types.Date{Time: r.Date},
		StartPage:   r.StartPage,
		EndPage:     r.EndPage,
		IsCompleted: r.IsCompleted,
		CreatedAt:   r.CreatedAt,
	}
	if r.ActualEndPage != nil {
		resp.ActualEndPage = r.ActualEndPage
	}
	if r.Memo != nil {
		resp.Memo = r.Memo
	}
	if r.CompletedAt != nil {
		resp.CompletedAt = r.CompletedAt
	}
	return resp
}

func toDailyTaskListResponse(r *usecase.DailyTaskListResult) oapi.DailyTaskListResponse {
	tasks := make([]oapi.DailyTaskResponse, len(r.Tasks))
	for i, t := range r.Tasks {
		tasks[i] = toDailyTaskResponse(t)
	}

	return oapi.DailyTaskListResponse{
		Date:  openapi_types.Date{Time: r.Date},
		Tasks: tasks,
		Summary: oapi.DailyTaskSummary{
			Total:      r.Total,
			Completed:  r.Completed,
			TotalPages: r.TotalPages,
		},
	}
}

func handleDailyTaskError(ctx echo.Context, err error) error {
	switch err {
	case usecase.ErrTaskNotFound:
		return ctx.JSON(http.StatusNotFound, oapi.ErrorResponse{Error: "task not found"})
	case usecase.ErrTaskAlreadyComplete:
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "task already completed"})
	case usecase.ErrTasksAlreadyExist:
		return ctx.JSON(http.StatusConflict, oapi.ErrorResponse{Error: "tasks already exist for this date"})
	case usecase.ErrInvalidTask:
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid task data"})
	case usecase.ErrInvalidDate:
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "date must be today or later"})
	default:
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}
}
