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

func (r *availabilityRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Availability, error) {
	rows, err := r.client.Availability.
		Query().
		Where(availability.HasUserWith(entUser.IDEQ(userID))).
		Order(ent.Asc(availability.FieldDayOfWeek)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Availability, len(rows))
	for i, row := range rows {
		result[i] = toAvailabilityModel(row, userID)
	}
	return result, nil
}

func (r *availabilityRepository) UpsertBatch(ctx context.Context, userID uuid.UUID, items []*model.Availability) ([]*model.Availability, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	// 既存レコードを全削除
	_, err = tx.Availability.
		Delete().
		Where(availability.HasUserWith(entUser.IDEQ(userID))).
		Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// 新規レコードを一括作成
	builders := make([]*ent.AvailabilityCreate, len(items))
	for i, item := range items {
		builders[i] = tx.Availability.
			Create().
			SetDayOfWeek(item.DayOfWeek).
			SetHours(item.Hours).
			SetUserID(userID)
	}

	created, err := tx.Availability.CreateBulk(builders...).Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	result := make([]*model.Availability, len(created))
	for i, row := range created {
		result[i] = toAvailabilityModel(row, userID)
	}
	return result, nil
}

func toAvailabilityModel(e *ent.Availability, userID uuid.UUID) *model.Availability {
	return &model.Availability{
		ID:        e.ID,
		UserID:    userID,
		DayOfWeek: e.DayOfWeek,
		Hours:     e.Hours,
	}
}
