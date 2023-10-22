package seeders

import (
	"github.com/OltLatifi/cv-builder-back/models"
	"gorm.io/gorm"
)

func SeedUserStatuses(db *gorm.DB) {
	statuses := []models.UserStatus{
		{Name: "Active", Description: "User is active and can fully interact with the platform."},
		{Name: "Inactive", Description: "User account is inactive and may need verification."},
		{Name: "Suspended", Description: "User has been temporarily banned from the platform due to violation of terms."},
		{Name: "Banned", Description: "User has been permanently banned from using the platform."},
		{Name: "Pending", Description: "User registration is pending approval or verification."},
		{Name: "Archived", Description: "User data is archived and the account is no longer active."},
	}

	for _, status := range statuses {
		var currentStatus models.UserStatus

		db.Where("name = ?", status.Name).First(&currentStatus)

		if currentStatus.ID == 0 {
			db.Create(&status)
		}
	}
}