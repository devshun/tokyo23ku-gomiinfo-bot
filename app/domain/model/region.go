package model

import "time"

type Region struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	WardID    int       `gorm:"not null" json:"ward_id,omitempty"`
	Name      string    `gorm:"size:255;not null" json:"name,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	Ward      Ward      `gorm:"foreignKey:WardID" json:"ward,omitempty"`
}
