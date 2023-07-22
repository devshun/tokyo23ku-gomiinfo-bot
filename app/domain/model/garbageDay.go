package model

import (
	"fmt"
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
	RegionID          int         `gorm:"not null" json:"regionId,omitempty"`
	GarbageType       GarbageType `gorm:"not null" json:"garbageType,omitempty"`
	DayOfWeek         Weekday     `gorm:"not null" json:"dayOfWeek,omitempty"`
	WeekNumberOfMonth int         `gorm:"" json:"weekNumberOfMonth,omitempty"`
	CreatedAt         time.Time   `gorm:"autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt         time.Time   `gorm:"autoUpdateTime" json:"updatedAt,omitempty"`
	Region            Region      `gorm:"foreignKey:RegionID" json:"region,omitempty"`
}

func (gd *GarbageDay) Format() string {
	dayOfWeek := gd.DayOfWeek.String()

	if gd.WeekNumberOfMonth > 0 {
		dayOfWeek = fmt.Sprintf("第%d%s", gd.WeekNumberOfMonth, dayOfWeek)
	}

	return dayOfWeek
}
