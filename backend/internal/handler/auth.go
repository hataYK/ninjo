package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/hatamotoyuki/ninjo/backend/internal/handler/dto"
	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

// AuthHandler は認証関連のHTTPハンドラ。
type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// Signup はユーザー新規登録。
// POST /api/v1/auth/signup
func (h *AuthHandler) Signup(c echo.Context) error {
	var req dto.SignupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	result, err := h.authUsecase.Signup(c.Request().Context(), usecase.SignupInput{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		if err == usecase.ErrEmailAlreadyExists {
			return c.JSON(http.StatusConflict, map[string]string{"error": "email already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	setTokenCookies(c, result.AccessToken, result.RefreshToken)

	return c.JSON(http.StatusCreated, dto.UserResponse{
		ID:          result.User.ID.String(),
		Email:       result.User.Email,
		DisplayName: result.User.DisplayName,
	})
}

// Login はログイン。
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	result, err := h.authUsecase.Login(c.Request().Context(), usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	setTokenCookies(c, result.AccessToken, result.RefreshToken)

	return c.JSON(http.StatusOK, dto.UserResponse{
		ID:          result.User.ID.String(),
		Email:       result.User.Email,
		DisplayName: result.User.DisplayName,
	})
}

// Logout はログアウト。Cookieを削除する。
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c echo.Context) error {
	clearTokenCookies(c)
	return c.NoContent(http.StatusNoContent)
}

// Refresh はアクセストークンを更新する。
// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "refresh token required"})
	}

	newAccessToken, err := h.authUsecase.Refresh(c.Request().Context(), cookie.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired refresh token"})
	}

	setAccessTokenCookie(c, newAccessToken)

	return c.JSON(http.StatusOK, map[string]string{"message": "token refreshed"})
}

// setTokenCookies はアクセストークンとリフレッシュトークンをCookieにセットする。
func setTokenCookies(c echo.Context, accessToken, refreshToken string) {
	setAccessTokenCookie(c, accessToken)

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth/refresh",
		HttpOnly: true,
		Secure:   false, // 開発環境ではfalse。本番ではtrue
		SameSite: http.SameSiteLaxMode,
		MaxAge:   30 * 24 * 60 * 60, // 30日
	})
}

func setAccessTokenCookie(c echo.Context, accessToken string) {
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // 開発環境ではfalse。本番ではtrue
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60, // 1時間
	})
}

// clearTokenCookies はCookieを削除する（MaxAge=-1で即時削除）。
func clearTokenCookies(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth/refresh",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

