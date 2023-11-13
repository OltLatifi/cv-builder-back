package models

import (
	"gorm.io/gorm"
)

type Experience struct {
	gorm.Model
	Name     string      `gorm:"type:varchar(255);not null" json:"name"`
	Company  string      `gorm:"type:varchar(255);not null" json:"company"`
	Position string      `gorm:"type:varchar(255);not null" json:"position"`
	StartAt  CustomDate  `json:"start_at"`
	EndAt    *CustomDate `json:"end_at"`
	UserId   uint        `json:"user_id"`
	User     User        `gorm:"foreignKey:UserId" json:"user"`
}
