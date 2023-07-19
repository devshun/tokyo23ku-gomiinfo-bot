package model

import "time"

type Ward struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Name      string    `gorm:"size:255;not null" json:"name,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	Regions   []Region  `gorm:"foreignKey:WardID" json:"regions,omitempty"`
}
