package model

import (
	"time"

	"github.com/google/uuid"
)

// User はユーザーのドメインモデル。
// ent の生成コードとは独立しており、ビジネスロジックで使う。
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	DisplayName  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
