package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/OltLatifi/cv-builder-back/helpers"
	"github.com/OltLatifi/cv-builder-back/initializers"
	"github.com/OltLatifi/cv-builder-back/models"
	"github.com/gin-gonic/gin"
)

func CreateExperience(c *gin.Context) {
	var body struct {
		Name     string  `json:"name" binding:"required"`
		Company  string  `json:"company" binding:"required"`
		Position string  `json:"position" binding:"required"`
		StartAt  string  `json:"start_at" binding:"required"`
		EndAt    *string `json:"end_at"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Validation failed", err)
		return
	}
	start_at, err := time.Parse("2006-01-02", body.StartAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Invalid date format",
			"description": "Expected format: YYYY-MM-DD",
		})
		return
	}

	var end_at *models.CustomDate
	if body.EndAt != nil {

		parsedEndAt, err := time.Parse("2006-01-02", *body.EndAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "Invalid date format",
				"description": "Expected format: YYYY-MM-DD",
			})
			return
		}
		end_at = &models.CustomDate{Time: parsedEndAt}
	}

	currentUser := c.MustGet("currentUser").(models.User)
	experience := models.Experience{
		Name:     body.Name,
		Company:  body.Company,
		Position: body.Position,
		StartAt:  models.CustomDate{Time: start_at},
		EndAt:    end_at,
		UserId:   currentUser.ID,
	}

	result := initializers.DB.Create(&experience)
	if result.Error != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Failed to create experience", result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"experience": experience,
	})
}

func GetExperiences(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	var experiences []models.Experience
	initializers.DB.Where("user_id = ?", currentUser.ID).Find(&experiences)

	c.JSON(http.StatusOK, gin.H{
		"experiences": experiences,
	})
}

func EditExperience(c *gin.Context) {
	var body struct {
		Name     string  `json:"name" binding:"required"`
		Company  string  `json:"company" binding:"required"`
		Position string  `json:"position" binding:"required"`
		StartAt  string  `json:"start_at" binding:"required"`
		EndAt    *string `json:"end_at"`
	}
	experienceID := c.Param("id")

	if err := c.ShouldBindJSON(&body); err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	var experience models.Experience
	result := initializers.DB.Where("id = ? AND user_id = ?", experienceID, currentUser.ID).First(&experience)
	if result.Error != nil {
		helpers.HandleError(c, http.StatusNotFound, "Experience not found", result.Error)
		return
	}

	start_at, err := time.Parse("2006-01-02", body.StartAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Invalid date format",
			"description": "Expected format: YYYY-MM-DD",
		})
		return
	}

	var end_at *models.CustomDate
	if body.EndAt != nil {

		parsedEndAt, err := time.Parse("2006-01-02", *body.EndAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "Invalid date format",
				"description": "Expected format: YYYY-MM-DD",
			})
			return
		}
		end_at = &models.CustomDate{Time: parsedEndAt}
	}

	experience.Name = body.Name
	experience.Company = body.Company
	experience.Position = body.Position
	experience.StartAt = models.CustomDate{Time: start_at}
	experience.EndAt = end_at

	result = initializers.DB.Save(&experience)
	if result.Error != nil {
		helpers.HandleError(c, http.StatusInternalServerError, "Failed to update experience", result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"experience": experience,
	})
}

func DeleteExperience(c *gin.Context) {
	experienceID := c.Param("id")
	currentUser := c.MustGet("currentUser").(models.User)

	result := initializers.DB.Where("id = ? AND user_id = ?", experienceID, currentUser.ID).Delete(&models.Experience{})
	if result.Error != nil {
		helpers.HandleError(c, http.StatusNotFound, "Experience not found", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		helpers.HandleError(c, http.StatusNotFound, "No experience found to delete", fmt.Errorf("no experience found to delete"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Experience deleted successfully",
	})
}
