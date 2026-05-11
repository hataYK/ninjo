package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/hatamotoyuki/ninjo/backend/internal/handler/middleware"
	"github.com/hatamotoyuki/ninjo/backend/internal/handler/oapi"
	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

// RegisterRoutes は OpenAPI から自動生成されたルーティングを登録する。
// 認証が必要なエンドポイントには OperationMiddlewares で JWT ミドルウェアを適用する。
func RegisterRoutes(e *echo.Echo, uc *usecase.Usecase) {
	h := NewHandler(uc)
	authMw := middleware.JWTAuth(uc.Auth())

	oapi.RegisterHandlersWithOptions(e, h, oapi.RegisterHandlersOptions{
		OperationMiddlewares: map[string][]echo.MiddlewareFunc{
			// 認証が必要なエンドポイント
			"getAvatar":          {authMw},
			"updateAvatar":       {authMw},
			"getAvailability":    {authMw},
			"updateAvailability": {authMw},
			"listPlans":          {authMw},
			"createPlan":         {authMw},
			"reviewPlan":         {authMw},
			"getPlan":            {authMw},
			"deletePlan":         {authMw},
			"extractSkills":      {authMw},
			"listSkills":         {authMw},
			"createSkill":        {authMw},
			"updateSkill":        {authMw},
			"deleteSkill":        {authMw},
		},
	})
}
