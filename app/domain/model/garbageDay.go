package model

import "time"

type GarbageDay struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	RegionID    int       `gorm:"not null"`
	GarbageType string    `gorm:"size:255;not null"`
	DayOfWeek   string    `gorm:"size:255;not null"`
	WeekOfMonth int       `gorm:""`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Region      Region    `gorm:"foreignKey:RegionID"`
}
