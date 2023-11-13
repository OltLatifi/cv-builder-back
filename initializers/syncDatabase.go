package initializers

import "github.com/OltLatifi/cv-builder-back/models"

func SyncDatabase() {
	DB.AutoMigrate(
		&models.Role{},
		&models.UserStatus{},
		&models.User{},
		&models.Language{},
		&models.VerificationToken{},
		&models.Experience{},
	)
}
