package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/repository"
)

var (
	ErrInvalidSkill    = errors.New("invalid skill data")
	ErrSkillNotFound   = errors.New("skill not found")
	ErrSkillDuplicate  = errors.New("skill name already exists")
)

// SkillUsecase はスキルに関するビジネスロジック。
type SkillUsecase struct {
	skillRepo repository.SkillRepository
}

func NewSkillUsecase(skillRepo repository.SkillRepository) *SkillUsecase {
	return &SkillUsecase{skillRepo: skillRepo}
}

// CreateSkillInput はスキル作成の入力。
type CreateSkillInput struct {
	Name     string
	Category string
	TaskID   *uuid.UUID
	Source   model.SkillSource
}

// SuggestedSkill はAI抽出されたスキル提案。
type SuggestedSkill struct {
	Name     string
	Category string
}

// ExtractSkills はタスクのメモからスキルを抽出する（保存しない）。
// TODO: 将来的にはClaude APIで抽出。現在は固定レスポンス。
func (uc *SkillUsecase) ExtractSkills(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) ([]SuggestedSkill, error) {
	// 将来: taskIDからメモを取得 → Claude APIでスキル抽出
	// 現在はテンプレートレスポンスを返す
	return []SuggestedSkill{
		{Name: "学習スキル", Category: "その他"},
	}, nil
}

// Create はスキルを作成する。
func (uc *SkillUsecase) Create(ctx context.Context, userID uuid.UUID, input CreateSkillInput) (*model.Skill, error) {
	if input.Name == "" || len(input.Name) > 200 {
		return nil, ErrInvalidSkill
	}
	if input.Category == "" || len(input.Category) > 100 {
		return nil, ErrInvalidSkill
	}

	skill := &model.Skill{
		UserID:   userID,
		TaskID:   input.TaskID,
		Name:     input.Name,
		Category: input.Category,
		Source:   input.Source,
	}

	created, err := uc.skillRepo.Create(ctx, skill)
	if err != nil {
		// ent のユニーク制約違反をチェック
		if isUniqueConstraintError(err) {
			return nil, ErrSkillDuplicate
		}
		return nil, err
	}
	return created, nil
}

// List はユーザーのスキル一覧を取得する。
func (uc *SkillUsecase) List(ctx context.Context, userID uuid.UUID, category *string) ([]*model.Skill, int, error) {
	var skills []*model.Skill
	var err error

	if category != nil && *category != "" {
		skills, err = uc.skillRepo.FindByUserIDAndCategory(ctx, userID, *category)
	} else {
		skills, err = uc.skillRepo.FindByUserID(ctx, userID)
	}
	if err != nil {
		return nil, 0, err
	}

	return skills, len(skills), nil
}

// Update はスキルを更新する。
func (uc *SkillUsecase) Update(ctx context.Context, userID uuid.UUID, skillID uuid.UUID, name, category string) (*model.Skill, error) {
	if name == "" || len(name) > 200 {
		return nil, ErrInvalidSkill
	}
	if category == "" || len(category) > 100 {
		return nil, ErrInvalidSkill
	}

	existing, err := uc.skillRepo.FindByID(ctx, skillID)
	if err != nil {
		return nil, err
	}
	if existing == nil || existing.UserID != userID {
		return nil, ErrSkillNotFound
	}

	existing.Name = name
	existing.Category = category

	updated, err := uc.skillRepo.Update(ctx, existing)
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil, ErrSkillDuplicate
		}
		return nil, err
	}
	return updated, nil
}

// Delete はスキルを削除する。
func (uc *SkillUsecase) Delete(ctx context.Context, userID uuid.UUID, skillID uuid.UUID) error {
	existing, err := uc.skillRepo.FindByID(ctx, skillID)
	if err != nil {
		return err
	}
	if existing == nil || existing.UserID != userID {
		return ErrSkillNotFound
	}

	return uc.skillRepo.Delete(ctx, skillID)
}

// isUniqueConstraintError は一意制約違反かどうかを判定する。
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "UNIQUE constraint failed")
}
