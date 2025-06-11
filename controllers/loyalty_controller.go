package controllers

import (
	"net/http"

	"github.com/gimhanr9/go-loyalty-api/services"
	"github.com/gin-gonic/gin"
)

// Initialize loyalty service
var loyaltyService = services.NewLoyaltyService()

type RedeemRequest struct {
	RewardTierID string `json:"reward_tier_id" binding:"required"`
}

type EarnRequest struct {
	Points int `json:"points" binding:"required"`
}

func RedeemPoints(c *gin.Context) {
	customerID := c.GetString("customer_id")

	var req RedeemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: reward_tier_id required"})
		return
	}

	if err := loyaltyService.RedeemPoints(customerID, req.RewardTierID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reward redeemed successfully"})
}

func EarnPoints(c *gin.Context) {
	customerID := c.GetString("customer_id")

	var req EarnRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Points <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: points must be > 0"})
		return
	}

	if err := loyaltyService.EarnPoints(customerID, req.Points); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Points earned successfully"})
}

func GetBalance(c *gin.Context) {
	customerID := c.GetString("customer_id")

	balance, err := loyaltyService.GetBalance(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func GetHistory(c *gin.Context) {
	customerID := c.GetString("customer_id")

	history, err := loyaltyService.GetHistory(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}
