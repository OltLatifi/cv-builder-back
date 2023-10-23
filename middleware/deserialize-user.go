package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/OltLatifi/cv-builder-back/initializers"
	"github.com/OltLatifi/cv-builder-back/models"
	"github.com/OltLatifi/cv-builder-back/utils"
	"github.com/gin-gonic/gin"
)

func DeserializeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var access_token string
		cookie, err := c.Cookie("access_token")

		authorizationHeader := c.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else if err == nil {
			access_token = cookie
		}

		if access_token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		sub, err := utils.ValidateToken(access_token, os.Getenv("ACCESS_TOKEN_PUBLIC_KEY"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		var user models.User
		result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}
