package model

import "time"

type Region struct {
	ID          int          `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	WardID      int          `gorm:"not null" json:"wardId,omitempty"`
	Name        string       `gorm:"size:255;not null" json:"name,omitempty"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updatedAt,omitempty"`
	Ward        Ward         `gorm:"foreignKey:WardID" json:"ward,omitempty"`
	GarbageDays []GarbageDay `gorm:"foreignKey:RegionID" json:"garbageDays,omitempty"`
}
