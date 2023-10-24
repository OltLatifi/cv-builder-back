package helpers

import "github.com/gin-gonic/gin"

func SetCookieWithEnvMaxAge(c *gin.Context, name, value, maxAgeKey string) {
	maxAge := GetEnvInt(maxAgeKey)
	c.SetCookie(name, value, maxAge*60, "/", "localhost", false, true)
}
