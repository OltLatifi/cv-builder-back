package models

import "gorm.io/gorm"

type UserStatus struct {
	gorm.Model
	Name        string `gorm:"size:255;not null"`
	Description string `gorm:"size:512"`
}