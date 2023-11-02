package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/OltLatifi/cv-builder-back/helpers"
	"github.com/OltLatifi/cv-builder-back/initializers"
	"github.com/OltLatifi/cv-builder-back/models"
	"github.com/OltLatifi/cv-builder-back/utils"
	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin/binding"
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
		StatusID:    5,
		RoleID:      body.RoleID,
	}

	result := initializers.DB.Create(&user)

	token := utils.CreateVerificationToken(user.ID)
	helpers.SendEmail(user.Email, token)

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
	type LoginBody struct {
		// Email    string `json:"email" binding:"required"`
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}

	var body LoginBody
	if err := c.ShouldBind(&body); err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	var user models.User
	initializers.DB.Preload("Status").First(&user, "email = ? OR username = ?", body.Identifier, body.Identifier)

	if user.ID == 0 {
		helpers.HandleError(c, http.StatusBadRequest, "Invalid email or password", fmt.Errorf("invalid email or password"))
		return
	}

	if user.StatusID != 1  {
		helpers.HandleError(c, http.StatusBadRequest, user.Status.Name, fmt.Errorf(user.Status.Description))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Invalid email or password", err)
		return
	}

	accessTokenDuration, err := helpers.GetEnvDuration("ACCESS_TOKEN_EXPIRED_IN")
	if err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Failed to get token duration", err)
		return
	}

	refreshTokenDuration, err := helpers.GetEnvDuration("REFRESH_TOKEN_EXPIRED_IN")
	if err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Failed to get refresh token duration", err)
		return
	}

	access_token, err := utils.CreateToken(accessTokenDuration, user.ID, os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Failed to create access token", err)
		return
	}

	refresh_token, err := utils.CreateToken(refreshTokenDuration, user.ID, os.Getenv("REFRESH_TOKEN_PRIVATE_KEY"))
	if err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Failed to create refresh token", err)
		return
	}

	helpers.SetCookieWithEnvMaxAge(c, "access_token", access_token, "ACCESS_TOKEN_MAXAGE")
	helpers.SetCookieWithEnvMaxAge(c, "refresh_token", refresh_token, "REFRESH_TOKEN_MAXAGE")
	c.SetCookie("logged_in", "true", helpers.GetEnvInt("ACCESS_TOKEN_MAXAGE")*60, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func RefreshToken(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		helpers.HandleError(c, http.StatusForbidden, "Could not refresh access token", err)
		return
	}

	sub, err := utils.ValidateToken(cookie, os.Getenv("REFRESH_TOKEN_PUBLIC_KEY"))
	if err != nil {
		helpers.HandleError(c, http.StatusForbidden, "Token validation failed", err)
		return
	}

	var user models.User
	if initializers.DB.First(&user, "ID = ?", fmt.Sprint(sub)).Error != nil {
		helpers.HandleError(c, http.StatusForbidden, "The user belonging to this token no longer exists", fmt.Errorf("User not found"))
		return
	}

	accessTokenDuration, err := helpers.GetEnvDuration("ACCESS_TOKEN_EXPIRED_IN")
	if err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Failed to get token duration", err)
		return
	}

	access_token, err := utils.CreateToken(accessTokenDuration, user.ID, os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		helpers.HandleError(c, http.StatusForbidden, "Failed to create access token", err)
		return
	}

	refreshTokenDuration, err := helpers.GetEnvDuration("REFRESH_TOKEN_EXPIRED_IN")
	if err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Failed to get refresh token duration", err)
		return
	}

	refresh_token, err := utils.CreateToken(refreshTokenDuration, user.ID, os.Getenv("REFRESH_TOKEN_PRIVATE_KEY"))
	if err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "Failed to create refresh token", err)
		return
	}

	helpers.SetCookieWithEnvMaxAge(c, "access_token", access_token, "ACCESS_TOKEN_MAXAGE")
	helpers.SetCookieWithEnvMaxAge(c, "refresh_token", refresh_token, "REFRESH_TOKEN_MAXAGE")
	c.SetCookie("logged_in", "true", helpers.GetEnvInt("ACCESS_TOKEN_MAXAGE")*60, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
