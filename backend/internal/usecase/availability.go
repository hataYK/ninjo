package usecase

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"

	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/repository"
)

var ErrInvalidAvailability = errors.New("invalid availability data")

// AvailabilityUsecase は可処分時間に関するビジネスロジック。
type AvailabilityUsecase struct {
	repo repository.AvailabilityRepository
}

func NewAvailabilityUsecase(repo repository.AvailabilityRepository) *AvailabilityUsecase {
	return &AvailabilityUsecase{repo: repo}
}

// AvailabilityItem はGET/PUTで使うDTO。
type AvailabilityItem struct {
	DayOfWeek int8
	Hours     float64
}

// AvailabilityResult はレスポンス用のDTO。
type AvailabilityResult struct {
	Items       []AvailabilityItem
	WeeklyTotal float64
}

// Get はユーザーの可処分時間を取得する。
// DBにレコードがない曜日はデフォルト0hとして補完する。
func (uc *AvailabilityUsecase) Get(ctx context.Context, userID uuid.UUID) (*AvailabilityResult, error) {
	rows, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 曜日→hoursのマップを作成
	hoursMap := make(map[int8]float64)
	for _, row := range rows {
		hoursMap[row.DayOfWeek] = row.Hours
	}

	// 全7曜日を補完
	items := make([]AvailabilityItem, 7)
	var total float64
	for i := int8(0); i <= 6; i++ {
		h := hoursMap[i]
		items[i] = AvailabilityItem{DayOfWeek: i, Hours: h}
		total += h
	}

	return &AvailabilityResult{Items: items, WeeklyTotal: total}, nil
}

// Update はユーザーの可処分時間を一括更新する。
func (uc *AvailabilityUsecase) Update(ctx context.Context, userID uuid.UUID, items []AvailabilityItem) (*AvailabilityResult, error) {
	if err := validateAvailabilityItems(items); err != nil {
		return nil, err
	}

	models := make([]*model.Availability, len(items))
	for i, item := range items {
		models[i] = &model.Availability{
			UserID:    userID,
			DayOfWeek: item.DayOfWeek,
			Hours:     item.Hours,
		}
	}

	_, err := uc.repo.UpsertBatch(ctx, userID, models)
	if err != nil {
		return nil, err
	}

	return uc.Get(ctx, userID)
}

func validateAvailabilityItems(items []AvailabilityItem) error {
	if len(items) != 7 {
		return ErrInvalidAvailability
	}

	seen := make(map[int8]bool)
	for _, item := range items {
		if item.DayOfWeek < 0 || item.DayOfWeek > 6 {
			return ErrInvalidAvailability
		}
		if item.Hours < 0 || item.Hours > 24 {
			return ErrInvalidAvailability
		}
		// 0.5刻みチェック
		if math.Mod(item.Hours, 0.5) != 0 {
			return ErrInvalidAvailability
		}
		if seen[item.DayOfWeek] {
			return ErrInvalidAvailability
		}
		seen[item.DayOfWeek] = true
	}

	return nil
}
