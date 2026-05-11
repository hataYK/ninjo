package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
)

// DailyTaskRepository はデイリータスクの永続化インターフェース。
type DailyTaskRepository interface {
	// Create はタスクを新規作成する。
	Create(ctx context.Context, task *model.DailyTask) (*model.DailyTask, error)
	// FindByID はIDでタスクを取得する。見つからなければnilを返す。
	FindByID(ctx context.Context, id uuid.UUID) (*model.DailyTask, error)
	// FindByUserAndDate はユーザーの指定日のタスク一覧を取得する。
	FindByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*model.DailyTask, error)
	// ExistsByPlanAndDate は指定計画・指定日のタスクが存在するか確認する。
	ExistsByPlanAndDate(ctx context.Context, planID uuid.UUID, date time.Time) (bool, error)
	// Update はタスクを更新する。
	Update(ctx context.Context, task *model.DailyTask) (*model.DailyTask, error)
	// FindLatestByPlanID は計画の最新の完了タスク（actual_end_page順）を取得する。
	FindLatestCompletedByPlanID(ctx context.Context, planID uuid.UUID) (*model.DailyTask, error)
}
