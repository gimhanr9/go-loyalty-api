package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gimhanr9/go-loyalty-api/controllers"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/earn", controllers.EarnPoints)
		api.POST("/redeem", controllers.RedeemPoints)
		api.GET("/balance", controllers.GetBalance)
		api.GET("/history", controllers.GetHistory)
	}
}
