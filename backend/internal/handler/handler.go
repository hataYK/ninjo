package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/hatamotoyuki/ninjo/backend/internal/handler/oapi"
	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

// Handler は oapi.ServerInterface を実装する。
// Usecase ファサードを持ち、各メソッドから必要なユースケースにアクセスする。
type Handler struct {
	uc *usecase.Usecase
}

func NewHandler(uc *usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}

// HealthCheck はヘルスチェック。
func (h *Handler) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, oapi.MessageResponse{Message: "ok"})
}

// コンパイル時に ServerInterface の実装を検証する。
var _ oapi.ServerInterface = (*Handler)(nil)
