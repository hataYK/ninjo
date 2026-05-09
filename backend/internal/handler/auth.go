package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/hatamotoyuki/ninjo/backend/internal/handler/oapi"
	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

// Signup はユーザー新規登録。
func (h *Handler) Signup(ctx echo.Context) error {
	var req oapi.SignupRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid request body"})
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: err.Error()})
	}

	result, err := h.uc.Auth().Signup(ctx.Request().Context(), usecase.SignupInput{
		Email:       string(req.Email),
		Password:    req.Password,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		if err == usecase.ErrEmailAlreadyExists {
			return ctx.JSON(http.StatusConflict, oapi.ErrorResponse{Error: "email already exists"})
		}
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	setTokenCookies(ctx, result.AccessToken, result.RefreshToken)

	return ctx.JSON(http.StatusCreated, oapi.UserResponse{
		Id:          result.User.ID,
		Email:       result.User.Email,
		DisplayName: result.User.DisplayName,
	})
}

// Login はログイン。
func (h *Handler) Login(ctx echo.Context) error {
	var req oapi.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: "invalid request body"})
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, oapi.ErrorResponse{Error: err.Error()})
	}

	result, err := h.uc.Auth().Login(ctx.Request().Context(), usecase.LoginInput{
		Email:    string(req.Email),
		Password: req.Password,
	})
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			return ctx.JSON(http.StatusUnauthorized, oapi.ErrorResponse{Error: "invalid email or password"})
		}
		return ctx.JSON(http.StatusInternalServerError, oapi.ErrorResponse{Error: "internal server error"})
	}

	setTokenCookies(ctx, result.AccessToken, result.RefreshToken)

	return ctx.JSON(http.StatusOK, oapi.UserResponse{
		Id:          result.User.ID,
		Email:       result.User.Email,
		DisplayName: result.User.DisplayName,
	})
}

// Logout はログアウト。
func (h *Handler) Logout(ctx echo.Context) error {
	clearTokenCookies(ctx)
	return ctx.NoContent(http.StatusNoContent)
}

// RefreshToken はアクセストークンを更新する。
func (h *Handler) RefreshToken(ctx echo.Context) error {
	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, oapi.ErrorResponse{Error: "refresh token required"})
	}

	newAccessToken, err := h.uc.Auth().Refresh(ctx.Request().Context(), cookie.Value)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, oapi.ErrorResponse{Error: "invalid or expired refresh token"})
	}

	setAccessTokenCookie(ctx, newAccessToken)

	return ctx.JSON(http.StatusOK, oapi.MessageResponse{Message: "token refreshed"})
}

// setTokenCookies はアクセストークンとリフレッシュトークンをCookieにセットする。
func setTokenCookies(c echo.Context, accessToken, refreshToken string) {
	setAccessTokenCookie(c, accessToken)

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth/refresh",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   30 * 24 * 60 * 60,
	})
}

func setAccessTokenCookie(c echo.Context, accessToken string) {
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60,
	})
}

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
