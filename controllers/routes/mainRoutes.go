package routes

import (
	"netto/controllers/gads"

	"github.com/gin-gonic/gin"
)

func InitGadsRoutes(router *gin.RouterGroup) {
	router.GET("/fetch-account", gads.FetchGoogle)

}
