package routes

import (
	"github.com/OltLatifi/cv-builder-back/controllers"
	"github.com/OltLatifi/cv-builder-back/middleware"
	"github.com/gin-gonic/gin"
)

func SetupApiRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// authentication
		v1.POST("/register", controllers.Register)
		v1.POST("/login", controllers.Login)
		v1.GET("/refresh", controllers.RefreshToken)
		v1.POST("/verify", controllers.VerifyUser)
		v1.POST("/logout", middleware.DeserializeUser(), controllers.Logout)

		// user
		v1.GET("/user-profile", middleware.DeserializeUser(), controllers.UserProfile)

		// Language
		languages := v1.Group("/languages")
		{
			languages.GET("/", middleware.DeserializeUser(), controllers.GetLanguages)
			languages.POST("/register", middleware.DeserializeUser(), controllers.SetUserLanguage)
			languages.PATCH("/update/:id", middleware.DeserializeUser(), controllers.EditUserLanguage)
			languages.DELETE("/delete/:id", middleware.DeserializeUser(), controllers.DeleteUserLanguage)
		}

		// Experience
		experience := v1.Group("/experience")
		{
			experience.GET("/", middleware.DeserializeUser(), controllers.GetExperiences)
			experience.POST("/register", middleware.DeserializeUser(), controllers.CreateExperience)
			experience.PUT("/update/:id", middleware.DeserializeUser(), controllers.EditExperience)
			experience.DELETE("/delete/:id", middleware.DeserializeUser(), controllers.DeleteExperience)
		}

	}

}
