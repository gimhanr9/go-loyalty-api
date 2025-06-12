package controllers

import (
	"net/http"

	"github.com/gimhanr9/go-loyalty-api/services"
	"github.com/gin-gonic/gin"
)

// Initialize loyalty service

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

	err := services.EarnPoints(customerID, req.Description, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	balance, err := services.GetBalance(customerID)
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

	balance, err := services.GetBalance(accountID)
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

	cursor := c.Query("cursor") // read from query param

	history, err := services.GetHistory(customerID, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(), // Return error message
		})
		return
	}

	c.JSON(http.StatusOK, history)
}

func GetRewardTiers(c *gin.Context) {
	accountID := c.Query("accountID")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing accountID"})
		return
	}

	percentage, err := services.GetDiscountPercentageByClosestRewardTier(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"discount_percentage": percentage})
}
