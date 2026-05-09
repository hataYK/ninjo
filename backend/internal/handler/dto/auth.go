package dto

// SignupRequest はサインアップのリクエストDTO。
type SignupRequest struct {
	Email       string `json:"email" validate:"required,email,max=255"`
	Password    string `json:"password" validate:"required,min=8"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
}

// LoginRequest はログインのリクエストDTO。
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse はユーザー情報のレスポンスDTO。
// パスワードハッシュは含めない。
type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}
