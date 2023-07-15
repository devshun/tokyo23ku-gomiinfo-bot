package model

import "time"

type Region struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	WardID    int       `gorm:"not null"`
	Name      string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Ward      Ward      `gorm:"foreignKey:WardID"`
}
