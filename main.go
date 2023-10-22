package main

import (
	"github.com/OltLatifi/cv-builder-back/initializers"
	"github.com/OltLatifi/cv-builder-back/routes/v1"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
	initializers.SyncDatabase()
	initializers.SyncSeeders()
}

func main() {
	router := gin.Default()
	routes.SetupApiRoutes(router)
	router.Run()
}