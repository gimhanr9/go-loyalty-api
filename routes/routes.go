package routes

import (
	"github.com/gimhanr9/go-loyalty-api/controllers"
	"github.com/gimhanr9/go-loyalty-api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Public
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)

	// Protected
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/earn", controllers.EarnPoints)
		protected.POST("/redeem", controllers.RedeemPoints)
		protected.GET("/balance", controllers.GetBalance)
		protected.GET("/history", controllers.GetHistory)
		protected.GET("/rewardtiers", controllers.GetRewardTiers)
	}
}
