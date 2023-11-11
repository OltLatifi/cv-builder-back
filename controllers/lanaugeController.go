package controllers

import (
	"fmt"
	"net/http"

	"github.com/OltLatifi/cv-builder-back/helpers"
	"github.com/OltLatifi/cv-builder-back/initializers"
	"github.com/OltLatifi/cv-builder-back/models"
	"github.com/gin-gonic/gin"
)

func SetUserLanguage(c *gin.Context) {
	var Body struct {
		Language string `json:"language" binding:"required"`
		Level    string `json:"level" binding:"required"`
	}

	if err := c.ShouldBind(&Body); err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	language := models.Language{Language: Body.Language, Level: Body.Level, UserId: currentUser.ID}

	result := initializers.DB.Create(&language)

	if result.Error != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Failed to create language", result.Error)
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

func EditUserLanguage(c *gin.Context) {
	var body struct {
		Language string `json:"language"`
		Level    string `json:"level"`
	}

	languageID := c.Param("id")

	if err := c.ShouldBindJSON(&body); err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)

	var language models.Language
	result := initializers.DB.Where("id = ? AND user_id = ?", languageID, currentUser.ID).First(&language)
	if result.Error != nil {
		helpers.HandleError(c, http.StatusNotFound, "Language not found", result.Error)
		return
	}

	language.Language = body.Language
	language.Level = body.Level

	result = initializers.DB.Save(&language)
	if result.Error != nil {
		helpers.HandleError(c, http.StatusInternalServerError, "Failed to update language", result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"language": language,
	})
}

func DeleteUserLanguage(c *gin.Context) {
	languageID := c.Param("id")

	currentUser := c.MustGet("currentUser").(models.User)

	var language models.Language
	result := initializers.DB.Where("id = ? AND user_id = ?", languageID, currentUser.ID).Delete(&language)
	if result.Error != nil {
		helpers.HandleError(c, http.StatusNotFound, "Language not found", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		helpers.HandleError(c, http.StatusNotFound, "No language found to delete", fmt.Errorf("no language found to delete"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Language deleted successfully",
	})
}
