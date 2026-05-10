package model

import "github.com/google/uuid"

// Availability は曜日ごとの可処分時間。
type Availability struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	DayOfWeek int8 // Go time.Weekday 準拠（0=日曜, 6=土曜）
	Hours     float64
}
