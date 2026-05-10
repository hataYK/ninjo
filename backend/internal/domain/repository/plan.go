package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
)

// PlanRepository は学習計画の永続化インターフェース。
type PlanRepository interface {
	// Create は計画を新規作成する。
	Create(ctx context.Context, plan *model.Plan) (*model.Plan, error)
	// FindByID はIDで計画を取得する。見つからなければnilを返す。
	FindByID(ctx context.Context, id uuid.UUID) (*model.Plan, error)
	// FindByUserID はユーザーの計画一覧を取得する。
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Plan, error)
	// CountActiveByUserID はユーザーのactive状態の計画数を返す。
	CountActiveByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	// Delete は計画を削除する。
	Delete(ctx context.Context, id uuid.UUID) error
}
