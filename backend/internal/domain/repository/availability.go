package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
)

// AvailabilityRepository は可処分時間の永続化インターフェース。
type AvailabilityRepository interface {
	// FindByUserID はユーザーの全曜日の可処分時間を取得する。
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Availability, error)
	// UpsertBatch はユーザーの可処分時間を一括upsertする。
	UpsertBatch(ctx context.Context, userID uuid.UUID, items []*model.Availability) ([]*model.Availability, error)
}
