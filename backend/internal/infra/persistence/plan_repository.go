package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/hatamotoyuki/ninjo/backend/ent"
	"github.com/hatamotoyuki/ninjo/backend/ent/plan"
	entUser "github.com/hatamotoyuki/ninjo/backend/ent/user"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/repository"
)

type planRepository struct {
	client *ent.Client
}

func NewPlanRepository(client *ent.Client) repository.PlanRepository {
	return &planRepository{client: client}
}

func (r *planRepository) Create(ctx context.Context, p *model.Plan) (*model.Plan, error) {
	created, err := r.client.Plan.
		Create().
		SetTitle(p.Title).
		SetTotalPages(p.TotalPages).
		SetStartDate(p.StartDate).
		SetTargetDate(p.TargetDate).
		SetStatus(plan.Status(p.Status)).
		SetNillableAiReview(p.AIReview).
		SetUserID(p.UserID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toPlanModel(created, p.UserID), nil
}

func (r *planRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Plan, error) {
	row, err := r.client.Plan.
		Query().
		Where(plan.IDEQ(id)).
		WithUser().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	userID := row.Edges.User.ID
	return toPlanModel(row, userID), nil
}

func (r *planRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Plan, error) {
	rows, err := r.client.Plan.
		Query().
		Where(plan.HasUserWith(entUser.IDEQ(userID))).
		Order(ent.Desc(plan.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	plans := make([]*model.Plan, len(rows))
	for i, row := range rows {
		plans[i] = toPlanModel(row, userID)
	}
	return plans, nil
}

func (r *planRepository) CountActiveByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	return r.client.Plan.
		Query().
		Where(
			plan.HasUserWith(entUser.IDEQ(userID)),
			plan.StatusEQ(plan.StatusActive),
		).
		Count(ctx)
}

func (r *planRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.Plan.
		DeleteOneID(id).
		Exec(ctx)
}

func toPlanModel(e *ent.Plan, userID uuid.UUID) *model.Plan {
	return &model.Plan{
		ID:         e.ID,
		UserID:     userID,
		Title:      e.Title,
		TotalPages: e.TotalPages,
		StartDate:  e.StartDate,
		TargetDate: e.TargetDate,
		Status:     model.PlanStatus(e.Status),
		AIReview:   e.AiReview,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}
