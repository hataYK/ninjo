package persistence

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/hatamotoyuki/ninjo/backend/ent"
	"github.com/hatamotoyuki/ninjo/backend/ent/dailytask"
	entPlan "github.com/hatamotoyuki/ninjo/backend/ent/plan"
	entUser "github.com/hatamotoyuki/ninjo/backend/ent/user"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/repository"
)

type dailyTaskRepository struct {
	client *ent.Client
}

func NewDailyTaskRepository(client *ent.Client) repository.DailyTaskRepository {
	return &dailyTaskRepository{client: client}
}

func (r *dailyTaskRepository) Create(ctx context.Context, task *model.DailyTask) (*model.DailyTask, error) {
	created, err := r.client.DailyTask.
		Create().
		SetDate(task.Date).
		SetStartPage(task.StartPage).
		SetEndPage(task.EndPage).
		SetPlanID(task.PlanID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toDailyTaskModel(created, task.PlanID), nil
}

func (r *dailyTaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.DailyTask, error) {
	row, err := r.client.DailyTask.
		Query().
		Where(dailytask.IDEQ(id)).
		WithPlan().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return toDailyTaskModel(row, row.Edges.Plan.ID), nil
}

func (r *dailyTaskRepository) FindByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*model.DailyTask, error) {
	rows, err := r.client.DailyTask.
		Query().
		Where(
			dailytask.DateEQ(date),
			dailytask.HasPlanWith(
				entPlan.HasUserWith(entUser.IDEQ(userID)),
			),
		).
		WithPlan().
		Order(ent.Asc(dailytask.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	tasks := make([]*model.DailyTask, len(rows))
	for i, row := range rows {
		tasks[i] = toDailyTaskModel(row, row.Edges.Plan.ID)
	}
	return tasks, nil
}

func (r *dailyTaskRepository) ExistsByPlanAndDate(ctx context.Context, planID uuid.UUID, date time.Time) (bool, error) {
	return r.client.DailyTask.
		Query().
		Where(
			dailytask.DateEQ(date),
			dailytask.HasPlanWith(entPlan.IDEQ(planID)),
		).
		Exist(ctx)
}

func (r *dailyTaskRepository) Update(ctx context.Context, task *model.DailyTask) (*model.DailyTask, error) {
	update := r.client.DailyTask.
		UpdateOneID(task.ID).
		SetStartPage(task.StartPage).
		SetEndPage(task.EndPage).
		SetIsCompleted(task.IsCompleted).
		SetNillableActualEndPage(task.ActualEndPage).
		SetNillableMemo(task.Memo).
		SetNillableCompletedAt(task.CompletedAt)

	updated, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}
	return toDailyTaskModel(updated, task.PlanID), nil
}

func (r *dailyTaskRepository) FindLatestCompletedByPlanID(ctx context.Context, planID uuid.UUID) (*model.DailyTask, error) {
	row, err := r.client.DailyTask.
		Query().
		Where(
			dailytask.HasPlanWith(entPlan.IDEQ(planID)),
			dailytask.IsCompletedEQ(true),
		).
		Order(ent.Desc(dailytask.FieldDate)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return toDailyTaskModel(row, planID), nil
}

func toDailyTaskModel(e *ent.DailyTask, planID uuid.UUID) *model.DailyTask {
	return &model.DailyTask{
		ID:            e.ID,
		PlanID:        planID,
		Date:          e.Date,
		StartPage:     e.StartPage,
		EndPage:       e.EndPage,
		ActualEndPage: e.ActualEndPage,
		IsCompleted:   e.IsCompleted,
		Memo:          e.Memo,
		CompletedAt:   e.CompletedAt,
		CreatedAt:     e.CreatedAt,
	}
}
