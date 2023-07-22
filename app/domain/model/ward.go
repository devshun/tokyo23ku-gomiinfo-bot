package model

import "time"

type Ward struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Name      string    `gorm:"size:255;not null" json:"name,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt,omitempty"`
	Regions   []Region  `gorm:"foreignKey:WardID" json:"regions,omitempty"`
}
