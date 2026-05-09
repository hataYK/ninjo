package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/hatamotoyuki/ninjo/backend/internal/handler/middleware"
	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

// RegisterRoutes はすべてのルーティングを登録する。
func RegisterRoutes(e *echo.Echo, authUsecase *usecase.AuthUsecase) {
	authHandler := NewAuthHandler(authUsecase)

	// 認証不要のエンドポイント
	auth := e.Group("/api/v1/auth")
	auth.POST("/signup", authHandler.Signup)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
	auth.POST("/refresh", authHandler.Refresh)

	// 認証が必要なエンドポイント（今後ここに追加）
	_ = e.Group("/api/v1", middleware.JWTAuth(authUsecase))
}
