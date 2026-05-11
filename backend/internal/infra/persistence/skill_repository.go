package persistence

import (
	"context"

	"github.com/google/uuid"

	"github.com/hatamotoyuki/ninjo/backend/ent"
	"github.com/hatamotoyuki/ninjo/backend/ent/skill"
	entUser "github.com/hatamotoyuki/ninjo/backend/ent/user"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/repository"
)

type skillRepository struct {
	client *ent.Client
}

func NewSkillRepository(client *ent.Client) repository.SkillRepository {
	return &skillRepository{client: client}
}

func (r *skillRepository) Create(ctx context.Context, s *model.Skill) (*model.Skill, error) {
	q := r.client.Skill.
		Create().
		SetName(s.Name).
		SetCategory(s.Category).
		SetSource(skill.Source(s.Source)).
		SetNillableTaskID(s.TaskID).
		SetUserID(s.UserID)

	created, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}
	return toSkillModel(created, s.UserID), nil
}

func (r *skillRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Skill, error) {
	row, err := r.client.Skill.
		Query().
		Where(skill.IDEQ(id)).
		WithUser().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	userID := row.Edges.User.ID
	return toSkillModel(row, userID), nil
}

func (r *skillRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Skill, error) {
	rows, err := r.client.Skill.
		Query().
		Where(skill.HasUserWith(entUser.IDEQ(userID))).
		Order(ent.Desc(skill.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	skills := make([]*model.Skill, len(rows))
	for i, row := range rows {
		skills[i] = toSkillModel(row, userID)
	}
	return skills, nil
}

func (r *skillRepository) FindByUserIDAndCategory(ctx context.Context, userID uuid.UUID, category string) ([]*model.Skill, error) {
	rows, err := r.client.Skill.
		Query().
		Where(
			skill.HasUserWith(entUser.IDEQ(userID)),
			skill.CategoryEQ(category),
		).
		Order(ent.Desc(skill.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	skills := make([]*model.Skill, len(rows))
	for i, row := range rows {
		skills[i] = toSkillModel(row, userID)
	}
	return skills, nil
}

func (r *skillRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	return r.client.Skill.
		Query().
		Where(skill.HasUserWith(entUser.IDEQ(userID))).
		Count(ctx)
}

func (r *skillRepository) Update(ctx context.Context, s *model.Skill) (*model.Skill, error) {
	updated, err := r.client.Skill.
		UpdateOneID(s.ID).
		SetName(s.Name).
		SetCategory(s.Category).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toSkillModel(updated, s.UserID), nil
}

func (r *skillRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.Skill.
		DeleteOneID(id).
		Exec(ctx)
}

func toSkillModel(e *ent.Skill, userID uuid.UUID) *model.Skill {
	return &model.Skill{
		ID:        e.ID,
		UserID:    userID,
		TaskID:    e.TaskID,
		Name:      e.Name,
		Category:  e.Category,
		Source:    model.SkillSource(e.Source),
		CreatedAt: e.CreatedAt,
	}
}
