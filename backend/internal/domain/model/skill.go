package model

import (
	"time"

	"github.com/google/uuid"
)

// SkillSource はスキルの作成元。
type SkillSource string

const (
	SkillSourceAI     SkillSource = "ai"
	SkillSourceManual SkillSource = "manual"
)

// Skill はスキルのドメインモデル。
type Skill struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TaskID    *uuid.UUID
	Name      string
	Category  string
	Source    SkillSource
	CreatedAt time.Time
}
