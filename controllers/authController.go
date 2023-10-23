package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/OltLatifi/cv-builder-back/helpers"
	"github.com/OltLatifi/cv-builder-back/initializers"
	"github.com/OltLatifi/cv-builder-back/models"
	"github.com/OltLatifi/cv-builder-back/utils"
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

func Login(c *gin.Context) {

	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Validation failed",
			"description": err.Error(),
		})
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"sub": user.ID,
	// 	"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	// })

	// tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Failed to create token",
	// 	})

	// 	return
	// }

	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	// c.JSON(http.StatusOK, gin.H{})

	ttl, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})

	}
	// Generate Tokens
	access_token, err := utils.CreateToken(ttl, user.ID, os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshTtl, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRED_IN"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})

	}

	refresh_token, err := utils.CreateToken(refreshTtl, user.ID, os.Getenv("REFRESH_TOKEN_PRIVATE_KEY"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	maxAgeStr := os.Getenv("ACCESS_TOKEN_MAXAGE")
	maxAge, err := strconv.Atoi(maxAgeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})

	}

	maxAgeStrRefreshToken := os.Getenv("REFRESH_TOKEN_MAXAGE")
	maxAgeRefreshToken, err := strconv.Atoi(maxAgeStrRefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})

	}

	c.SetCookie("access_token", access_token, maxAge*60, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refresh_token, maxAgeRefreshToken*60, "/", "localhost", false, true)
	c.SetCookie("logged_in", "true", maxAge*60, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func RefreshToken(c *gin.Context) {
	message := "could not refresh access token"

	cookie, err := c.Cookie("refresh_token")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	sub, err := utils.ValidateToken(cookie, os.Getenv("REFRESH_TOKEN_PUBLIC_KEY"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "ID = ?", fmt.Sprint(sub))

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	ttl, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})

	}

	access_token, err := utils.CreateToken(ttl, user.ID, os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	maxAgeStr := os.Getenv("ACCESS_TOKEN_MAXAGE")
	maxAge, err := strconv.Atoi(maxAgeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})

	}

	c.SetCookie("access_token", access_token, maxAge*60, "/", "localhost", false, true)
	c.SetCookie("logged_in", "true", maxAge*60, "/", "localhost", false, false)
	// If we want to prolong the refresh token we add one more line here
	// c.SetCookie("refresh_token", refresh_token, maxAgeRefreshToken*60, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
