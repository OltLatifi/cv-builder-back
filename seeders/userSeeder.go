package seeders

import (
	"time"

	"github.com/OltLatifi/cv-builder-back/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	var adminUser models.User
	db.Where("email = ?", "admin@admin.com").First(&adminUser)

	// If the admin user doesn't exist, create it
	if adminUser.ID == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		adminUser := models.User{
			Username:    "admin",
			Email:       "admin@admin.com",
			Password:    string(hashedPassword),
			FirstName:   "Admin",
			LastName:    "User",
			DateOfBirth: models.CustomDate{Time: time.Now()},
			Address:     "123 Admin St, Admin City",
			Bio:         "I am the admin user.",
			StatusID:    1,
			RoleID:      1,
		}
		db.Create(&adminUser)
	}
}
