package controllers

import (
	"mime/multipart"
	"net/http"
	"time"

	"github.com/OltLatifi/cv-builder-back/helpers"
	"github.com/OltLatifi/cv-builder-back/initializers"
	"github.com/OltLatifi/cv-builder-back/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// If there the statusId is not provided then we will set to 5 - pending
func defaultStatusID(statusID uint) uint {
	if statusID == 0 {
		return 5
	}
	return statusID
}

func Register(c *gin.Context) {
	var body struct {
		FirstName   string                `form:"first_name" binding:"required"`
		LastName    string                `form:"last_name" binding:"required"`
		Username    string                `form:"username" binding:"required"`
		Email       string                `form:"email" binding:"required,email"`
		Password    string                `form:"password" binding:"required"`
		DateOfBirth string                `form:"date_of_birth"`
		Image       *multipart.FileHeader `form:"image" binding:"omitempty"`
		Address     string                `form:"address"`
		Bio         string                `form:"bio"`
		StatusID    uint                  `form:"status_id" binding:"required"`
		RoleID      uint                  `form:"role_id" binding:"required"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Validation failed",
			"description": err.Error(),
		})
		return
	}

	var role models.Role
	if err := initializers.DB.First(&role, body.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Invalid Role ID",
			"description": err.Error(),
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Failed to hash password",
			"description": err.Error(),
		})
		return
	}
	dob, err := time.Parse("2006-01-02", body.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Invalid date format",
			"description": "Expected format: YYYY-MM-DD",
		})
		return
	}
	// Truncate the time from the parsed date
	// TODO: FIX THIS - STORES DATE TIME SHOULD STORE ONLY DATE
	dob = dob.Truncate(24 * time.Hour)

	// Image upload
	imagePath, err := helpers.UploadImage(c, "image", body.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Failed to upload image",
			"description": err.Error(),
		})
		return
	}

	user := models.User{
		Username:    body.Username,
		Email:       body.Email,
		Password:    string(hash),
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		DateOfBirth: models.CustomDate{Time: dob},
		Image:       imagePath,
		Address:     body.Address,
		Bio:         body.Bio,
		StatusID:    defaultStatusID(body.StatusID),
		RoleID:      body.RoleID,
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Failed to create user",
			"description": result.Error.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
