package controllers

import (
	"net/http"

	"github.com/OltLatifi/cv-builder-back/initializers"
	"github.com/OltLatifi/cv-builder-back/models"
	"github.com/gin-gonic/gin"
)

func SetUserLanguage(c *gin.Context) {
	var Body struct {
		Language string `json:"language" binding:"required"`
		Level	 string `json:"level" binding:"required"`
	}

	if err := c.ShouldBind(&Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Validation failed",
			"description": err.Error(),
		})
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	language := models.Language{Language: Body.Language, Level: Body.Level, UserId: currentUser.ID}

	result := initializers.DB.Create(&language)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Failed to create language",
			"description": result.Error.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"language": language,
	})
}

func GetLanguages(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	var languages []models.Language

	initializers.DB.Where("user_id = ?", currentUser.ID).Find(&languages)

	c.JSON(http.StatusOK, gin.H{
		"languages": languages,
	})
}