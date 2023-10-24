package helpers

import "github.com/gin-gonic/gin"

func HandleError(c *gin.Context, status int, message string, description error) {
	c.JSON(status, gin.H{
		"error":       message,
		"description": description.Error(),
	})
}
