package routes

import (
	"github.com/OltLatifi/cv-builder-back/controllers"
	"github.com/gin-gonic/gin"
)

func SetupApiRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", controllers.Register)
	}

}