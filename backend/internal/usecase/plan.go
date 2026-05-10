package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/hatamotoyuki/ninjo/backend/internal/domain/model"
	"github.com/hatamotoyuki/ninjo/backend/internal/domain/repository"
)

var (
	ErrInvalidPlan      = errors.New("invalid plan data")
	ErrPlanNotFound     = errors.New("plan not found")
	ErrTooManyPlans     = errors.New("too many active plans")
	ErrTargetDatePast   = errors.New("target date must be in the future")
)

const maxActivePlans = 10

// PlanUsecase は計画に関するビジネスロジック。
type PlanUsecase struct {
	planRepo         repository.PlanRepository
	availabilityRepo repository.AvailabilityRepository
}

func NewPlanUsecase(planRepo repository.PlanRepository, availabilityRepo repository.AvailabilityRepository) *PlanUsecase {
	return &PlanUsecase{planRepo: planRepo, availabilityRepo: availabilityRepo}
}

// PlanInput は計画作成・レビュー用の入力。
type PlanInput struct {
	Title      string
	TotalPages int
	TargetDate time.Time
}

// PlanReviewResult はAIレビューの結果。
type PlanReviewResult struct {
	DailyPages    float64
	TotalDays     int
	AvailableDays int
	ReviewMessage string
}

// PlanResult は計画のレスポンス用DTO。
type PlanResult struct {
	ID           uuid.UUID
	Title        string
	TotalPages   int
	StartDate    time.Time
	TargetDate   time.Time
	Status       model.PlanStatus
	AIReview     *string
	ProgressRate float64
	CreatedAt    time.Time
}

// PlanDetailResult は計画詳細のレスポンス用DTO。
type PlanDetailResult struct {
	PlanResult
	DailyPagesNeeded float64
	DaysRemaining    int
}

// Review は計画の入力内容をレビューする（保存しない）。
func (uc *PlanUsecase) Review(ctx context.Context, userID uuid.UUID, input PlanInput) (*PlanReviewResult, error) {
	if err := validatePlanInput(input); err != nil {
		return nil, err
	}

	today := time.Now().Truncate(24 * time.Hour)
	totalDays := int(input.TargetDate.Sub(today).Hours() / 24)

	// 可処分時間から勉強可能日を算出
	availableDays := uc.calcAvailableDays(ctx, userID, today, input.TargetDate)

	var dailyPages float64
	if availableDays > 0 {
		dailyPages = float64(input.TotalPages) / float64(availableDays)
	} else {
		// 可処分時間未設定の場合はtotalDaysで割る
		dailyPages = float64(input.TotalPages) / float64(totalDays)
	}
	dailyPages = math.Round(dailyPages*10) / 10

	reviewMessage := generateReviewMessage(input.Title, input.TotalPages, totalDays, availableDays, dailyPages)

	return &PlanReviewResult{
		DailyPages:    dailyPages,
		TotalDays:     totalDays,
		AvailableDays: availableDays,
		ReviewMessage: reviewMessage,
	}, nil
}

// Create は計画を作成する。
func (uc *PlanUsecase) Create(ctx context.Context, userID uuid.UUID, input PlanInput) (*PlanResult, error) {
	if err := validatePlanInput(input); err != nil {
		return nil, err
	}

	// active計画数チェック
	count, err := uc.planRepo.CountActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= maxActivePlans {
		return nil, ErrTooManyPlans
	}

	// AIレビューを生成して保存
	review, err := uc.Review(ctx, userID, input)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	plan := &model.Plan{
		UserID:     userID,
		Title:      input.Title,
		TotalPages: input.TotalPages,
		StartDate:  now.Truncate(24 * time.Hour),
		TargetDate: input.TargetDate,
		Status:     model.PlanStatusActive,
		AIReview:   &review.ReviewMessage,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	saved, err := uc.planRepo.Create(ctx, plan)
	if err != nil {
		return nil, err
	}

	return &PlanResult{
		ID:           saved.ID,
		Title:        saved.Title,
		TotalPages:   saved.TotalPages,
		StartDate:    saved.StartDate,
		TargetDate:   saved.TargetDate,
		Status:       saved.Status,
		AIReview:     saved.AIReview,
		ProgressRate: 0.0,
		CreatedAt:    saved.CreatedAt,
	}, nil
}

// List はユーザーの計画一覧を取得する。
func (uc *PlanUsecase) List(ctx context.Context, userID uuid.UUID) ([]PlanResult, error) {
	plans, err := uc.planRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	results := make([]PlanResult, len(plans))
	for i, p := range plans {
		results[i] = PlanResult{
			ID:           p.ID,
			Title:        p.Title,
			TotalPages:   p.TotalPages,
			StartDate:    p.StartDate,
			TargetDate:   p.TargetDate,
			Status:       p.Status,
			AIReview:     p.AIReview,
			ProgressRate: 0.0, // TODO: daily_tasksから導出
			CreatedAt:    p.CreatedAt,
		}
	}

	return results, nil
}

// Get はIDで計画詳細を取得する。
func (uc *PlanUsecase) Get(ctx context.Context, userID uuid.UUID, planID uuid.UUID) (*PlanDetailResult, error) {
	plan, err := uc.planRepo.FindByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	if plan == nil || plan.UserID != userID {
		return nil, ErrPlanNotFound
	}

	today := time.Now().Truncate(24 * time.Hour)
	daysRemaining := int(plan.TargetDate.Sub(today).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	// 残り可能日で1日あたりの必要ページ数を算出
	availableDays := uc.calcAvailableDays(ctx, userID, today, plan.TargetDate)
	var dailyPagesNeeded float64
	if availableDays > 0 {
		dailyPagesNeeded = float64(plan.TotalPages) / float64(availableDays) // TODO: 進捗分を引く
	} else if daysRemaining > 0 {
		dailyPagesNeeded = float64(plan.TotalPages) / float64(daysRemaining)
	}
	dailyPagesNeeded = math.Round(dailyPagesNeeded*10) / 10

	return &PlanDetailResult{
		PlanResult: PlanResult{
			ID:           plan.ID,
			Title:        plan.Title,
			TotalPages:   plan.TotalPages,
			StartDate:    plan.StartDate,
			TargetDate:   plan.TargetDate,
			Status:       plan.Status,
			AIReview:     plan.AIReview,
			ProgressRate: 0.0, // TODO: daily_tasksから導出
			CreatedAt:    plan.CreatedAt,
		},
		DailyPagesNeeded: dailyPagesNeeded,
		DaysRemaining:    daysRemaining,
	}, nil
}

// Delete は計画を削除する。
func (uc *PlanUsecase) Delete(ctx context.Context, userID uuid.UUID, planID uuid.UUID) error {
	plan, err := uc.planRepo.FindByID(ctx, planID)
	if err != nil {
		return err
	}
	if plan == nil || plan.UserID != userID {
		return ErrPlanNotFound
	}

	return uc.planRepo.Delete(ctx, planID)
}

// calcAvailableDays は期間内で可処分時間 > 0 の日数を算出する。
func (uc *PlanUsecase) calcAvailableDays(ctx context.Context, userID uuid.UUID, from, to time.Time) int {
	avail, err := uc.availabilityRepo.FindByUserID(ctx, userID)
	if err != nil || avail == nil {
		return 0
	}

	count := 0
	for d := from; d.Before(to); d = d.AddDate(0, 0, 1) {
		dayOfWeek := int(d.Weekday())
		if avail.Hours(dayOfWeek) > 0 {
			count++
		}
	}
	return count
}

func validatePlanInput(input PlanInput) error {
	if input.Title == "" || len(input.Title) > 200 {
		return ErrInvalidPlan
	}
	if input.TotalPages < 1 || input.TotalPages > 10000 {
		return ErrInvalidPlan
	}
	today := time.Now().Truncate(24 * time.Hour)
	if !input.TargetDate.After(today) {
		return ErrTargetDatePast
	}
	return nil
}

func generateReviewMessage(title string, totalPages, totalDays, availableDays int, dailyPages float64) string {
	// TODO: 将来的にはClaude APIで生成。現在は簡易テンプレート。
	if availableDays > 0 {
		return fmt.Sprintf(
			"残り%d日、勉強できる日が%d日あるので、1日あたり約%.0fページ必要です。コツコツ進めていきましょう！",
			totalDays, availableDays, math.Ceil(dailyPages),
		)
	}
	return fmt.Sprintf(
		"残り%d日で%dページ、1日あたり約%.0fページ必要です。可処分時間を設定すると、より正確なレビューができますよ！",
		totalDays, totalPages, math.Ceil(dailyPages),
	)
}
