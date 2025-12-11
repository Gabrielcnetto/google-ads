package main

import (
	"netto/controllers/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))
	api := router.Group("/api")
	routes.InitGadsRoutes(api.Group("/google-ads"))
	router.Run(":8080")
}
