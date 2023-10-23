package controllers

import (
	"net/http"

	"github.com/OltLatifi/cv-builder-back/models"
	"github.com/gin-gonic/gin"
)

func UserProfile(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": currentUser}})

}
