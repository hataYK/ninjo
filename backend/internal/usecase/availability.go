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

// AvailabilityInput はPUT用の入力。
type AvailabilityInput struct {
	DayOfWeek int8
	Hours     float64
}

// AvailabilityResult はレスポンス用のDTO。
type AvailabilityResult struct {
	Items       []AvailabilityInput
	WeeklyTotal float64
}

// Get はユーザーの可処分時間を取得する。
// レコードがなければデフォルト全0hを返す。
func (uc *AvailabilityUsecase) Get(ctx context.Context, userID uuid.UUID) (*AvailabilityResult, error) {
	avail, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if avail == nil {
		// デフォルト: 全曜日0h
		avail = &model.Availability{UserID: userID}
	}

	return toAvailabilityResult(avail), nil
}

// Update はユーザーの可処分時間を一括更新する。
func (uc *AvailabilityUsecase) Update(ctx context.Context, userID uuid.UUID, items []AvailabilityInput) (*AvailabilityResult, error) {
	if err := validateAvailabilityItems(items); err != nil {
		return nil, err
	}

	// 入力をモデルに変換
	avail := &model.Availability{UserID: userID}
	for _, item := range items {
		switch item.DayOfWeek {
		case 0:
			avail.SunHours = item.Hours
		case 1:
			avail.MonHours = item.Hours
		case 2:
			avail.TueHours = item.Hours
		case 3:
			avail.WedHours = item.Hours
		case 4:
			avail.ThuHours = item.Hours
		case 5:
			avail.FriHours = item.Hours
		case 6:
			avail.SatHours = item.Hours
		}
	}

	saved, err := uc.repo.Upsert(ctx, avail)
	if err != nil {
		return nil, err
	}

	return toAvailabilityResult(saved), nil
}

func toAvailabilityResult(avail *model.Availability) *AvailabilityResult {
	items := make([]AvailabilityInput, 7)
	for i := int8(0); i <= 6; i++ {
		items[i] = AvailabilityInput{
			DayOfWeek: i,
			Hours:     avail.Hours(int(i)),
		}
	}
	return &AvailabilityResult{
		Items:       items,
		WeeklyTotal: avail.WeeklyTotal(),
	}
}

func validateAvailabilityItems(items []AvailabilityInput) error {
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
