package model

import "time"

type Ward struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
