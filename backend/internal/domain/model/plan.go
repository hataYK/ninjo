package model

import (
	"time"

	"github.com/google/uuid"
)

// PlanStatus は計画の状態。
type PlanStatus string

const (
	PlanStatusActive    PlanStatus = "active"
	PlanStatusCompleted PlanStatus = "completed"
	PlanStatusPaused    PlanStatus = "paused"
)

// Plan は学習計画のドメインモデル。
type Plan struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Title      string
	TotalPages int
	StartDate  time.Time
	TargetDate time.Time
	Status     PlanStatus
	AIReview   *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
