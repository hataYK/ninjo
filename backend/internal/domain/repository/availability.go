package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
)

// AvailabilityRepository は可処分時間の永続化インターフェース。
type AvailabilityRepository interface {
	// FindByUserID はユーザーの可処分時間を取得する。レコードがなければnilを返す。
	FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Availability, error)
	// Upsert はユーザーの可処分時間を作成または更新する。
	Upsert(ctx context.Context, avail *model.Availability) (*model.Availability, error)
}
