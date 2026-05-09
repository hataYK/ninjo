package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/hatamotoyuki/ninjo/backend/internal/handler/oapi"
	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

// RegisterRoutes は OpenAPI から自動生成されたルーティングを登録する。
func RegisterRoutes(e *echo.Echo, uc *usecase.Usecase) {
	h := NewHandler(uc)
	oapi.RegisterHandlers(e, h)
}
