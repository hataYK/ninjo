package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
)

// UserRepository はユーザーの永続化インターフェース。
// domain 層で定義し、infra 層で実装する（依存性逆転）。
type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}
