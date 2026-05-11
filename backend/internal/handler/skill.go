package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
	"github.com/hatamotoyuki/ninjo/backend/internal/handler/middleware"
	"github.com/hatamotoyuki/ninjo/backend/internal/handler/oapi"
	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

// ExtractSkills はタスクのメモからスキルを抽出する（保存しない）。
func (h *Handler) ExtractSkills(ctx echo.Context, taskId openapi_types.UUID) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	suggested, err := h.uc.Skill().ExtractSkills(ctx.Request().Context(), userID, taskId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	items := make([]oapi.SuggestedSkill, len(suggested))
	for i, s := range suggested {
		items[i] = oapi.SuggestedSkill{
			Name:     s.Name,
			Category: s.Category,
		}
	}

	return ctx.JSON(http.StatusOK, oapi.ExtractSkillsResponse{SuggestedSkills: items})
}

// ListSkills はユーザーの全スキル一覧を取得する。
func (h *Handler) ListSkills(ctx echo.Context, params oapi.ListSkillsParams) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	skills, totalCount, err := h.uc.Skill().List(ctx.Request().Context(), userID, params.Category)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	items := make([]oapi.SkillResponse, len(skills))
	for i, s := range skills {
		items[i] = toSkillResponse(s)
	}

	return ctx.JSON(http.StatusOK, oapi.SkillListResponse{
		Skills:     items,
		TotalCount: totalCount,
	})
}

// CreateSkill はスキルを手動追加する。
func (h *Handler) CreateSkill(ctx echo.Context) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	var req oapi.CreateSkillRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid request body"})
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: err.Error()})
	}

	input := usecase.CreateSkillInput{
		Name:     req.Name,
		Category: req.Category,
		Source:   model.SkillSourceManual,
	}
	if req.TaskId != nil {
		taskID := uuid.UUID(*req.TaskId)
		input.TaskID = &taskID
	}

	created, err := h.uc.Skill().Create(ctx.Request().Context(), userID, input)
	if err != nil {
		return handleSkillError(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, toSkillResponse(created))
}

// UpdateSkill はスキルの名前やカテゴリを更新する。
func (h *Handler) UpdateSkill(ctx echo.Context, skillId openapi_types.UUID) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	var req oapi.UpdateSkillRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid request body"})
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: err.Error()})
	}

	updated, err := h.uc.Skill().Update(ctx.Request().Context(), userID, skillId, req.Name, req.Category)
	if err != nil {
		return handleSkillError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, toSkillResponse(updated))
}

// DeleteSkill はスキルを削除する。
func (h *Handler) DeleteSkill(ctx echo.Context, skillId openapi_types.UUID) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	err := h.uc.Skill().Delete(ctx.Request().Context(), userID, skillId)
	if err != nil {
		if err == usecase.ErrSkillNotFound {
			return ctx.JSON(http.StatusNotFound, oapi.ErrorResponse{Error: "skill not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	return ctx.NoContent(http.StatusNoContent)
}

func toSkillResponse(s *model.Skill) oapi.SkillResponse {
	resp := oapi.SkillResponse{
		Id:        s.ID,
		Name:      s.Name,
		Category:  s.Category,
		Source:    oapi.SkillResponseSource(s.Source),
		CreatedAt: s.CreatedAt,
	}
	if s.TaskID != nil {
		taskID := openapi_types.UUID(*s.TaskID)
		resp.TaskId = &taskID
	}
	return resp
}

func handleSkillError(ctx echo.Context, err error) error {
	switch err {
	case usecase.ErrInvalidSkill:
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid skill data"})
	case usecase.ErrSkillNotFound:
		return ctx.JSON(http.StatusNotFound, oapi.ErrorResponse{Error: "skill not found"})
	case usecase.ErrSkillDuplicate:
		return ctx.JSON(http.StatusConflict, oapi.ErrorResponse{Error: "skill name already exists"})
	default:
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}
}
