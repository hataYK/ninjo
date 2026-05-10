package model

import "github.com/google/uuid"

// Availability は週の曜日ごとの可処分時間。1ユーザーにつき1つ。
type Availability struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	SunHours float64
	MonHours float64
	TueHours float64
	WedHours float64
	ThuHours float64
	FriHours float64
	SatHours float64
}

// WeeklyTotal は週の合計可処分時間を返す。
func (a *Availability) WeeklyTotal() float64 {
	return a.SunHours + a.MonHours + a.TueHours + a.WedHours + a.ThuHours + a.FriHours + a.SatHours
}

// Hours は曜日（time.Weekday準拠: 0=日〜6=土）の可処分時間を返す。
func (a *Availability) Hours(dayOfWeek int) float64 {
	switch dayOfWeek {
	case 0:
		return a.SunHours
	case 1:
		return a.MonHours
	case 2:
		return a.TueHours
	case 3:
		return a.WedHours
	case 4:
		return a.ThuHours
	case 5:
		return a.FriHours
	case 6:
		return a.SatHours
	default:
		return 0
	}
}
