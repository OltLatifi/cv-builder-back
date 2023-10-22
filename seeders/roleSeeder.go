package seeders

import (
	"github.com/OltLatifi/cv-builder-back/models"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	var roleAdmin models.Role
	var roleUser models.Role

	db.Where("name = ?", "Admin").First(&roleAdmin)
	db.Where("name = ?", "User").First(&roleUser)

	if roleAdmin.ID == 0 {
		roleAdmin = models.Role{Name: "Admin"}
		db.Create(&roleAdmin)
	}

	if roleUser.ID == 0 {
		roleUser = models.Role{Name: "User"}
		db.Create(&roleUser)
	}
}