package usecase

import (
	"github.com/hatamotoyuki/ninjo/backend/internal/infra"
)

// UsecaseConfig はユースケース層が必要とする全ての依存。
// 外部依存（DataStore, API Key等）はここに集約する。
type UsecaseConfig struct {
	DS        *infra.DataStore
	JWTSecret string
}

// Usecase はビジネスロジック層のファサード。
// 全ユースケースへのアクセスを一元管理する。
// 機能追加時はここにメソッドを足すだけでよい。
type Usecase struct {
	config UsecaseConfig
}

func NewUsecase(config UsecaseConfig) *Usecase {
	return &Usecase{config: config}
}

func (u *Usecase) Auth() *AuthUsecase {
	return NewAuthUsecase(u.config.DS.User(), u.config.JWTSecret)
}

// 今後追加:
// func (u *Usecase) Plan() *PlanUsecase { ... }
// func (u *Usecase) DailyTask() *DailyTaskUsecase { ... }
// func (u *Usecase) Availability() *AvailabilityUsecase { ... }
