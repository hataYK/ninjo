package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/hatamotoyuki/ninjo/backend/internal/handler/middleware"
	"github.com/hatamotoyuki/ninjo/backend/internal/handler/oapi"
	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

// GetAvatar はアバター設定とスキルサマリーを取得する。
func (h *Handler) GetAvatar(ctx echo.Context) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	result, err := h.uc.Avatar().Get(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	categories := make([]oapi.SkillCategoryCount, len(result.SkillCategories))
	for i, sc := range result.SkillCategories {
		categories[i] = oapi.SkillCategoryCount{
			Category: sc.Category,
			Count:    sc.Count,
		}
	}

	return ctx.JSON(http.StatusOK, oapi.AvatarResponse{
		AvatarPresetId:  result.AvatarPresetID,
		SkillCount:      result.SkillCount,
		SkillCategories: categories,
	})
}

// UpdateAvatar はアバターのプリセットを設定/変更する。
func (h *Handler) UpdateAvatar(ctx echo.Context) error {
	userID := ctx.Get(middleware.ContextKeyUserID).(uuid.UUID)

	var req oapi.UpdateAvatarRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid request body"})
	}

	presetID, err := h.uc.Avatar().Update(ctx.Request().Context(), userID, req.AvatarPresetId)
	if err != nil {
		if err == usecase.ErrInvalidAvatarPreset {
			return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid avatar preset id"})
		}
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	return ctx.JSON(http.StatusOK, oapi.AvatarPresetResponse{
		AvatarPresetId: presetID,
	})
}
