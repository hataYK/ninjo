package model

import (
	"time"

	"github.com/google/uuid"
)

// DailyTask はデイリータスクのドメインモデル。
type DailyTask struct {
	ID            uuid.UUID
	PlanID        uuid.UUID
	Date          time.Time
	StartPage     int
	EndPage       int
	ActualEndPage *int
	IsCompleted   bool
	Memo          *string
	CompletedAt   *time.Time
	CreatedAt     time.Time
}
