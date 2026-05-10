package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/hatamotoyuki/ninjo/backend/ent"
	"github.com/hatamotoyuki/ninjo/backend/ent/availability"
	entUser "github.com/hatamotoyuki/ninjo/backend/ent/user"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/repository"
)

type availabilityRepository struct {
	client *ent.Client
}

func NewAvailabilityRepository(client *ent.Client) repository.AvailabilityRepository {
	return &availabilityRepository{client: client}
}

func (r *availabilityRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*model.Availability, error) {
	row, err := r.client.Availability.
		Query().
		Where(availability.HasUserWith(entUser.IDEQ(userID))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return toAvailabilityModel(row, userID), nil
}

func (r *availabilityRepository) Upsert(ctx context.Context, avail *model.Availability) (*model.Availability, error) {
	// 既存レコードを検索
	existing, err := r.client.Availability.
		Query().
		Where(availability.HasUserWith(entUser.IDEQ(avail.UserID))).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}

	if existing != nil {
		// 更新
		updated, err := existing.Update().
			SetSunHours(avail.SunHours).
			SetMonHours(avail.MonHours).
			SetTueHours(avail.TueHours).
			SetWedHours(avail.WedHours).
			SetThuHours(avail.ThuHours).
			SetFriHours(avail.FriHours).
			SetSatHours(avail.SatHours).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		return toAvailabilityModel(updated, avail.UserID), nil
	}

	// 新規作成
	created, err := r.client.Availability.
		Create().
		SetSunHours(avail.SunHours).
		SetMonHours(avail.MonHours).
		SetTueHours(avail.TueHours).
		SetWedHours(avail.WedHours).
		SetThuHours(avail.ThuHours).
		SetFriHours(avail.FriHours).
		SetSatHours(avail.SatHours).
		SetUserID(avail.UserID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toAvailabilityModel(created, avail.UserID), nil
}

func toAvailabilityModel(e *ent.Availability, userID uuid.UUID) *model.Availability {
	return &model.Availability{
		ID:       e.ID,
		UserID:   userID,
		SunHours: e.SunHours,
		MonHours: e.MonHours,
		TueHours: e.TueHours,
		WedHours: e.WedHours,
		ThuHours: e.ThuHours,
		FriHours: e.FriHours,
		SatHours: e.SatHours,
	}
}
