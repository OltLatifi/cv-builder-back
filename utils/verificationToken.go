package utils

import (
    "log"
    "time"
    "github.com/google/uuid"
    "github.com/OltLatifi/cv-builder-back/models"
    "github.com/OltLatifi/cv-builder-back/helpers"
    "github.com/OltLatifi/cv-builder-back/initializers"
)

func GenerateUUID() (string, error) {
    u, err := uuid.NewRandom()
    if err != nil {
        return "", err
    }
    return u.String(), nil
}

func CreateVerificationToken(user_id uint) string {
    u, err := GenerateUUID()
    if err != nil {
        log.Println("Error while generating UUID: ", err)
        return ""
    }

    fastForwardBy := helpers.GetEnvInt("VERIFICATION_EMAIL_AGE")

    currentDateTime := time.Now()
    futureDate := currentDateTime.Add(time.Duration(fastForwardBy) * time.Minute)

    token := models.VerificationToken{
        Token:      u, 
        UserID:     user_id,
        ValidUntil: futureDate,
    }

    result := initializers.DB.Create(&token)
    if result.Error != nil {
        log.Println("Error creating verification token: ", result.Error)
    }

	return u
}
