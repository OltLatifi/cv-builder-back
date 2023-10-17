package main

import (
	"github.com/OltLatifi/cv-builder-back/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	router := gin.Default()

	// router.POST("/signup", controllers.SignUp)
	// router.POST("/login", controllers.Login)
	// router.GET("/validate", middleware.RequireAuth, controllers.Validate)

	router.Run()
}