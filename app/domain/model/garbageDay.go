package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Weekday int

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
	ID                int       `gorm:"primaryKey;autoIncrement"`
	RegionID          int       `gorm:"not null"`
	GarbageType       string    `gorm:"size:255;not null"`
	DayOfWeek         Weekday   `gorm:""`
	WeekNumberOfMonth int       `gorm:""`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
	Region            Region    `gorm:"foreignKey:RegionID"`
}

func (w Weekday) String() string {
	return weekdayNames[w]
}

func FindWeekday(s string) (Weekday, int, error) {
	for i, name := range weekdayNames {
		// 曜日を取得
		if strings.Contains(s, name) {
			// 第何週目かを取得
			re := regexp.MustCompile(`第(\d)`)

			match := re.FindStringSubmatch(s)

			if len(match) > 0 {
				weekNum, err := strconv.Atoi(match[1])

				if err != nil {
					return 0, 0, err
				}

				return Weekday(i), weekNum, nil
			}

			return Weekday(i), 0, nil
		}
	}
	return 0, 0, fmt.Errorf("invalid: %s", s)
}
