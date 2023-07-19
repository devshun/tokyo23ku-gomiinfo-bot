package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type GarbageDay struct {
	ID                int         `gorm:"primaryKey;autoIncrement"`
	RegionID          int         `gorm:"not null"`
	GarbageType       GarbageType `gorm:"not null"`
	DayOfWeek         Weekday     `gorm:"not null"`
	WeekNumberOfMonth int         `gorm:""`
	CreatedAt         time.Time   `gorm:"autoCreateTime"`
	UpdatedAt         time.Time   `gorm:"autoUpdateTime"`
	Region            Region      `gorm:"foreignKey:RegionID"`
}

type Weekday int

const (
	Sunday Weekday = iota + 1
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

var weekdayMap = map[Weekday]string{
	Sunday:    "日曜日",
	Monday:    "月曜日",
	Tuesday:   "火曜日",
	Wednesday: "水曜日",
	Thursday:  "木曜日",
	Friday:    "金曜日",
	Saturday:  "土曜日",
}

func (w Weekday) String() string {
	return weekdayMap[w]
}

func FindWeekday(s string) (Weekday, int, error) {
	for k, v := range weekdayMap {
		// 曜日を取得
		if strings.Contains(s, v) {
			// 第何週目かを取得
			re := regexp.MustCompile(`第(\d)`)

			match := re.FindStringSubmatch(s)

			if len(match) > 0 {
				weekNum, err := strconv.Atoi(match[1])

				if err != nil {
					return 0, 0, err
				}

				return Weekday(k), weekNum, nil
			}

			return Weekday(k), 0, nil
		}
	}
	return 0, 0, fmt.Errorf("invalid: %s", s)
}

type GarbageType int

const (
	Burnable GarbageType = iota + 1
	NonBurnable
	Recyclable
)

var GarbageTypeMap = map[GarbageType]string{
	Burnable:    "燃えるゴミ",
	NonBurnable: "燃えないゴミ",
	Recyclable:  "資源ゴミ",
}

func (g GarbageType) String() string {
	return GarbageTypeMap[g]
}
