package usecase

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/repository"
)

var (
	ErrTaskNotFound        = errors.New("task not found")
	ErrTaskAlreadyComplete = errors.New("task already completed")
	ErrTasksAlreadyExist   = errors.New("tasks already exist for this date")
	ErrInvalidTask         = errors.New("invalid task data")
	ErrInvalidDate         = errors.New("invalid date")
)

// DailyTaskUsecase はデイリータスクに関するビジネスロジック。
type DailyTaskUsecase struct {
	taskRepo         repository.DailyTaskRepository
	planRepo         repository.PlanRepository
	availabilityRepo repository.AvailabilityRepository
}

func NewDailyTaskUsecase(
	taskRepo repository.DailyTaskRepository,
	planRepo repository.PlanRepository,
	availabilityRepo repository.AvailabilityRepository,
) *DailyTaskUsecase {
	return &DailyTaskUsecase{
		taskRepo:         taskRepo,
		planRepo:         planRepo,
		availabilityRepo: availabilityRepo,
	}
}

// DailyTaskResult はタスクのレスポンス用DTO。
type DailyTaskResult struct {
	ID            uuid.UUID
	PlanID        uuid.UUID
	PlanTitle     string
	Date          time.Time
	StartPage     int
	EndPage       int
	ActualEndPage *int
	IsCompleted   bool
	Memo          *string
	CompletedAt   *time.Time
	CreatedAt     time.Time
}

// DailyTaskListResult はタスク一覧のレスポンス用DTO。
type DailyTaskListResult struct {
	Date       time.Time
	Tasks      []DailyTaskResult
	Total      int
	Completed  int
	TotalPages int
}

// List は指定日のデイリータスクを取得する。
func (uc *DailyTaskUsecase) List(ctx context.Context, userID uuid.UUID, date time.Time) (*DailyTaskListResult, error) {
	tasks, err := uc.taskRepo.FindByUserAndDate(ctx, userID, date)
	if err != nil {
		return nil, err
	}

	// plan情報を取得してresultに変換
	results := make([]DailyTaskResult, len(tasks))
	completed := 0
	totalPages := 0
	for i, t := range tasks {
		plan, err := uc.planRepo.FindByID(ctx, t.PlanID)
		if err != nil {
			return nil, err
		}
		planTitle := ""
		if plan != nil {
			planTitle = plan.Title
		}
		results[i] = toDailyTaskResult(t, planTitle)
		if t.IsCompleted {
			completed++
		}
		totalPages += t.EndPage - t.StartPage + 1
	}

	return &DailyTaskListResult{
		Date:       date,
		Tasks:      results,
		Total:      len(results),
		Completed:  completed,
		TotalPages: totalPages,
	}, nil
}

// Generate は指定日のデイリータスクを自動生成する。
func (uc *DailyTaskUsecase) Generate(ctx context.Context, userID uuid.UUID, date time.Time) (*DailyTaskListResult, error) {
	today := time.Now().Truncate(24 * time.Hour)
	if date.Before(today) {
		return nil, ErrInvalidDate
	}

	// ユーザーのactive計画を取得
	plans, err := uc.planRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 可処分時間を取得
	avail, _ := uc.availabilityRepo.FindByUserID(ctx, userID)

	var createdTasks []*model.DailyTask
	hasConflict := false

	for _, plan := range plans {
		if plan.Status != model.PlanStatusActive {
			continue
		}

		// 既にタスクがあるか確認
		exists, err := uc.taskRepo.ExistsByPlanAndDate(ctx, plan.ID, date)
		if err != nil {
			return nil, err
		}
		if exists {
			hasConflict = true
			continue
		}

		// 現在の進捗ページを取得
		currentPage := 0
		latest, err := uc.taskRepo.FindLatestCompletedByPlanID(ctx, plan.ID)
		if err != nil {
			return nil, err
		}
		if latest != nil && latest.ActualEndPage != nil {
			currentPage = *latest.ActualEndPage
		}

		remainPages := plan.TotalPages - currentPage
		if remainPages <= 0 {
			continue
		}

		// 今日のページ配分を計算
		todayPages := uc.calcTodayPages(avail, plan, date, currentPage)

		startPage := currentPage + 1
		endPage := startPage + todayPages - 1
		if endPage > plan.TotalPages {
			endPage = plan.TotalPages
		}

		task := &model.DailyTask{
			PlanID:    plan.ID,
			Date:      date,
			StartPage: startPage,
			EndPage:   endPage,
		}

		created, err := uc.taskRepo.Create(ctx, task)
		if err != nil {
			return nil, err
		}
		createdTasks = append(createdTasks, created)
	}

	// 全計画が既にタスク生成済みだった場合
	if hasConflict && len(createdTasks) == 0 {
		return nil, ErrTasksAlreadyExist
	}

	// 結果を組み立て
	results := make([]DailyTaskResult, len(createdTasks))
	totalPages := 0
	for i, t := range createdTasks {
		plan, _ := uc.planRepo.FindByID(ctx, t.PlanID)
		planTitle := ""
		if plan != nil {
			planTitle = plan.Title
		}
		results[i] = toDailyTaskResult(t, planTitle)
		totalPages += t.EndPage - t.StartPage + 1
	}

	return &DailyTaskListResult{
		Date:       date,
		Tasks:      results,
		Total:      len(results),
		Completed:  0,
		TotalPages: totalPages,
	}, nil
}

// Update はタスクのページ範囲を更新する。
func (uc *DailyTaskUsecase) Update(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, startPage, endPage int) (*DailyTaskResult, error) {
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	// 所有権チェック
	plan, err := uc.planRepo.FindByID(ctx, task.PlanID)
	if err != nil {
		return nil, err
	}
	if plan == nil || plan.UserID != userID {
		return nil, ErrTaskNotFound
	}

	if task.IsCompleted {
		return nil, ErrTaskAlreadyComplete
	}

	// バリデーション
	if startPage < 1 || endPage < startPage || endPage > plan.TotalPages {
		return nil, ErrInvalidTask
	}

	task.StartPage = startPage
	task.EndPage = endPage

	updated, err := uc.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	result := toDailyTaskResult(updated, plan.Title)
	return &result, nil
}

// Complete はタスクを完了にする。
func (uc *DailyTaskUsecase) Complete(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, actualEndPage int, memo *string) (*DailyTaskResult, error) {
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	// 所有権チェック
	plan, err := uc.planRepo.FindByID(ctx, task.PlanID)
	if err != nil {
		return nil, err
	}
	if plan == nil || plan.UserID != userID {
		return nil, ErrTaskNotFound
	}

	if task.IsCompleted {
		return nil, ErrTaskAlreadyComplete
	}

	// バリデーション
	if actualEndPage < 1 || actualEndPage > plan.TotalPages {
		return nil, ErrInvalidTask
	}
	if memo != nil && len(*memo) > 5000 {
		return nil, ErrInvalidTask
	}

	now := time.Now()
	task.ActualEndPage = &actualEndPage
	task.Memo = memo
	task.IsCompleted = true
	task.CompletedAt = &now

	updated, err := uc.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	// 計画完了チェック: actual_end_page が total_pages に達したら completed にする
	if actualEndPage >= plan.TotalPages {
		plan.Status = model.PlanStatusCompleted
		plan.UpdatedAt = now
		// PlanRepository に UpdateStatus メソッドがないので、ここでは省略
		// TODO: plan のステータス更新を追加
	}

	result := toDailyTaskResult(updated, plan.Title)
	return &result, nil
}

// calcTodayPages は今日分のページ配分を計算する。
func (uc *DailyTaskUsecase) calcTodayPages(avail *model.Availability, plan *model.Plan, date time.Time, currentPage int) int {
	remainPages := plan.TotalPages - currentPage
	remainDays := int(plan.TargetDate.Sub(date).Hours()/24) + 1
	if remainDays <= 0 {
		return remainPages
	}

	if avail == nil || avail.WeeklyTotal() == 0 {
		// 均等割り
		return int(math.Ceil(float64(remainPages) / float64(remainDays)))
	}

	// 残り日の可処分時間合計を計算
	totalRemainHours := 0.0
	for d := date; !d.After(plan.TargetDate); d = d.AddDate(0, 0, 1) {
		totalRemainHours += avail.Hours(int(d.Weekday()))
	}

	if totalRemainHours == 0 {
		return int(math.Ceil(float64(remainPages) / float64(remainDays)))
	}

	todayHours := avail.Hours(int(date.Weekday()))
	todayPages := float64(remainPages) * (todayHours / totalRemainHours)
	result := int(math.Ceil(todayPages))
	if result < 1 && remainPages > 0 {
		result = 1
	}
	return result
}

func toDailyTaskResult(t *model.DailyTask, planTitle string) DailyTaskResult {
	return DailyTaskResult{
		ID:            t.ID,
		PlanID:        t.PlanID,
		PlanTitle:     planTitle,
		Date:          t.Date,
		StartPage:     t.StartPage,
		EndPage:       t.EndPage,
		ActualEndPage: t.ActualEndPage,
		IsCompleted:   t.IsCompleted,
		Memo:          t.Memo,
		CompletedAt:   t.CompletedAt,
		CreatedAt:     t.CreatedAt,
	}
}
