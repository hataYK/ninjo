package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/hatamotoyuki/ninjo/backend/internal/domain/repository"
)

const (
	DefaultAvatarPresetID = "preset_01"
	MaxAvatarPresetID     = 12
)

var ErrInvalidAvatarPreset = errors.New("invalid avatar preset id")

// AvatarUsecase はアバターに関するビジネスロジック。
type AvatarUsecase struct {
	userRepo repository.UserRepository
}

func NewAvatarUsecase(userRepo repository.UserRepository) *AvatarUsecase {
	return &AvatarUsecase{userRepo: userRepo}
}

// SkillCategoryCount はスキルカテゴリごとの件数。
type SkillCategoryCount struct {
	Category string
	Count    int
}

// AvatarResult はGETレスポンス用のDTO。
type AvatarResult struct {
	AvatarPresetID  string
	SkillCount      int
	SkillCategories []SkillCategoryCount
}

// Get はユーザーのアバター情報を取得する。
func (uc *AvatarUsecase) Get(ctx context.Context, userID uuid.UUID) (*AvatarResult, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	presetID := DefaultAvatarPresetID
	if user.AvatarPresetID != nil {
		presetID = *user.AvatarPresetID
	}

	// F9(Skills)が未実装のため、スキルは0件で返す
	return &AvatarResult{
		AvatarPresetID:  presetID,
		SkillCount:      0,
		SkillCategories: []SkillCategoryCount{},
	}, nil
}

// Update はアバターのプリセットIDを更新する。
func (uc *AvatarUsecase) Update(ctx context.Context, userID uuid.UUID, presetID string) (string, error) {
	if !isValidPresetID(presetID) {
		return "", ErrInvalidAvatarPreset
	}

	user, err := uc.userRepo.UpdateAvatarPresetID(ctx, userID, presetID)
	if err != nil {
		return "", err
	}

	if user.AvatarPresetID != nil {
		return *user.AvatarPresetID, nil
	}
	return DefaultAvatarPresetID, nil
}

func isValidPresetID(id string) bool {
	for i := 1; i <= MaxAvatarPresetID; i++ {
		if id == fmt.Sprintf("preset_%02d", i) {
			return true
		}
	}
	return false
}
