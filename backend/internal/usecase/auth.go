package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/repository"
)

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidToken        = errors.New("invalid or expired token")
)

// AuthUsecase は認証に関するビジネスロジック。
type AuthUsecase struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// SignupInput はサインアップの入力。
type SignupInput struct {
	Email       string
	Password    string
	DisplayName string
}

// AuthResult は認証成功時の結果。
type AuthResult struct {
	User         *model.User
	AccessToken  string
	RefreshToken string
}

// Signup はユーザーを新規登録する。
func (uc *AuthUsecase) Signup(ctx context.Context, input SignupInput) (*AuthResult, error) {
	// メールアドレスの重複チェック
	existing, _ := uc.userRepo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// パスワードをbcryptでハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: string(hash),
		DisplayName:  input.DisplayName,
	}

	created, err := uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// JWTトークン生成
	accessToken, err := uc.generateToken(created.ID, 1*time.Hour)
	if err != nil {
		return nil, err
	}
	refreshToken, err := uc.generateToken(created.ID, 30*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         created,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// LoginInput はログインの入力。
type LoginInput struct {
	Email    string
	Password string
}

// Login はメール+パスワードで認証する。
func (uc *AuthUsecase) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// bcryptでパスワードを検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := uc.generateToken(user.ID, 1*time.Hour)
	if err != nil {
		return nil, err
	}
	refreshToken, err := uc.generateToken(user.ID, 30*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Refresh はリフレッシュトークンから新しいアクセストークンを生成する。
func (uc *AuthUsecase) Refresh(ctx context.Context, refreshToken string) (string, error) {
	userID, err := uc.ValidateToken(refreshToken)
	if err != nil {
		return "", ErrInvalidToken
	}

	// ユーザーが存在するか確認
	_, err = uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", ErrInvalidToken
	}

	return uc.generateToken(userID, 1*time.Hour)
}

// ValidateToken はJWTトークンを検証し、user_idを返す。
func (uc *AuthUsecase) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return uc.jwtSecret, nil
	})
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	return uuid.Parse(sub)
}

// generateToken はJWTトークンを生成する。
// ペイロードに user_id (sub) と有効期限 (exp) を含む。
func (uc *AuthUsecase) generateToken(userID uuid.UUID, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(uc.jwtSecret)
}
