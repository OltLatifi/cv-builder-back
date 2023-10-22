package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName   string     `gorm:"type:varchar(100);not null" form:"first_name"`
	LastName    string     `gorm:"type:varchar(100);not null" form:"last_name"`
	Username    string     `gorm:"type:varchar(100);unique;not null" form:"username"`
	Email       string     `gorm:"type:varchar(255);unique;not null" form:"email"`
	Password    string     `gorm:"type:varchar(255);not null" form:"password"`
	DateOfBirth CustomDate `form:"date_of_birth"`
	Image       string     `gorm:"type:text" form:"image"`
	Address     string     `gorm:"type:text" form:"address"`
	Bio         string     `gorm:"type:text" form:"bio"`
	StatusID    uint       `form:"status_id"`
	Status      UserStatus `gorm:"foreignKey:StatusID" form:"status"`
	RoleID      uint       `form:"role_id"`
	Role        Role       `gorm:"foreignKey:RoleID" form:"role"`
}
