package model

import (
	"fmt"
	"strings"
	"time"
)

type Weekday int

// [堤通一丁目 土曜日 火曜日 金曜日 第2月曜日 第4月曜日]

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

var weekdayNames = [...]string{
	"日曜日",
	"月曜日",
	"火曜日",
	"水曜日",
	"木曜日",
	"金曜日",
	"土曜日",
}

type GarbageDay struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	RegionID    int       `gorm:"not null"`
	GarbageType string    `gorm:"size:255;not null"`
	DayOfWeek   Weekday   `gorm:""`
	WeekOfMonth int       `gorm:""`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Region      Region    `gorm:"foreignKey:RegionID"`
}

func (w Weekday) String() string {
	return weekdayNames[w]
}

func FindWeekday(s string) (Weekday, error) {
	for i, name := range weekdayNames {
		if strings.Contains(s, name) {
			return Weekday(i), nil
		}
	}
	return -1, fmt.Errorf("invalid: %s", s)
}
