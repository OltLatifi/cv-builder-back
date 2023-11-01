package models

import "gorm.io/gorm"

type Language struct {
	gorm.Model
	Language string `gorm:"type:varchar(255);not null"  json:"language"`
	Level    string `gorm:"type:varchar(255);not null" json:"level"`
	UserId   uint   `json:"user_id"`
	User     User   `gorm:"foreignKey:UserId" json:"user"`
}