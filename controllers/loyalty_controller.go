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

// func RedeemPoints(c *gin.Context) {
// 	customerID := c.GetString("customer_id")

// 	var req RedeemRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: reward_tier_id required"})
// 		return
// 	}

// 	if err := loyaltyService.RedeemPoints(customerID, req.RewardTierID); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Reward redeemed successfully"})
// }

func EarnPoints(c *gin.Context) {
	customerID := c.GetString("customer_id")

	var req struct {
		Amount      int    `json:"amount"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 || req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: amount must be > 0 and description is required"})
		return
	}

	err := loyaltyService.EarnPoints(customerID, req.Description, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	balance, err := loyaltyService.GetBalance(customerID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Points earned successfully, but failed to fetch balance",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Points earned successfully",
		"balance": balance,
	})
}
func GetBalance(c *gin.Context) {
	accountID := c.GetString("customer_id")

	balance, err := loyaltyService.GetBalance(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func GetHistory(c *gin.Context) {
	customerID := c.GetString("customer_id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing customer_id in context"})
		return
	}

	history, err := loyaltyService.GetHistory(customerID)
	if err != nil {
		// Optional: log the error for internal monitoring
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve loyalty history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}
