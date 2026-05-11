package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
)

// SkillRepository はスキルの永続化インターフェース。
type SkillRepository interface {
	// Create はスキルを新規作成する。
	Create(ctx context.Context, skill *model.Skill) (*model.Skill, error)
	// FindByID はIDでスキルを取得する。見つからなければnilを返す。
	FindByID(ctx context.Context, id uuid.UUID) (*model.Skill, error)
	// FindByUserID はユーザーの全スキルを取得する。
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Skill, error)
	// FindByUserIDAndCategory はカテゴリでフィルタしたスキル一覧を取得する。
	FindByUserIDAndCategory(ctx context.Context, userID uuid.UUID, category string) ([]*model.Skill, error)
	// CountByUserID はユーザーのスキル数を返す。
	CountByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	// Update はスキルを更新する。
	Update(ctx context.Context, skill *model.Skill) (*model.Skill, error)
	// Delete はスキルを削除する。
	Delete(ctx context.Context, id uuid.UUID) error
}
