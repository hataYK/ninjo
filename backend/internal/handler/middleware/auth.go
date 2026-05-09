package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

// ContextKeyUserID はコンテキストに格納する user_id のキー。
const ContextKeyUserID = "user_id"

// JWTAuth は認証ミドルウェア。
// Cookie から access_token を取得し、JWT を検証する。
// 認証成功時は user_id をコンテキストにセットする。
func JWTAuth(authUsecase *usecase.AuthUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("access_token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "authentication required",
				})
			}

			userID, err := authUsecase.ValidateToken(cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}

			c.Set(ContextKeyUserID, userID)
			return next(c)
		}
	}
}
