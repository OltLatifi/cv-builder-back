package models

import (
    "time"
    "gorm.io/gorm"
)

type VerificationToken struct {
    gorm.Model
    Token      string    `gorm:"size:512;not null" json:"token"`
    UserID     uint      `json:"user_id"`
    ValidUntil time.Time `json:"valid_until"`
}
