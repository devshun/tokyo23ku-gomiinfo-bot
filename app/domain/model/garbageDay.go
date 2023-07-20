package model

import (
	"time"
)

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

var WeekdayMap = map[Weekday]string{
	Sunday:    "日曜日",
	Monday:    "月曜日",
	Tuesday:   "火曜日",
	Wednesday: "水曜日",
	Thursday:  "木曜日",
	Friday:    "金曜日",
	Saturday:  "土曜日",
}

func (w Weekday) String() string {
	return WeekdayMap[w]
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

type GarbageDay struct {
	ID                int         `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	RegionID          int         `gorm:"not null" json:"region_id,omitempty"`
	GarbageType       GarbageType `gorm:"not null" json:"garbage_type,omitempty"`
	DayOfWeek         Weekday     `gorm:"not null" json:"day_of_week,omitempty"`
	WeekNumberOfMonth int         `gorm:"" json:"week_number_of_month,omitempty"`
	CreatedAt         time.Time   `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt         time.Time   `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	Region            Region      `gorm:"foreignKey:RegionID" json:"region,omitempty"`
}
