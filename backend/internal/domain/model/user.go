package model

import (
	"time"

	"github.com/google/uuid"
)

// User はユーザーのドメインモデル。
// ent の生成コードとは独立しており、ビジネスロジックで使う。
type User struct {
	ID             uuid.UUID
	Email          string
	PasswordHash   string
	DisplayName    string
	AvatarPresetID *string // nil の場合はデフォルト（"preset_01"）
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
