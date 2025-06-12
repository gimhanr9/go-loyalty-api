package controllers

import (
	"net/http"

	"github.com/gimhanr9/go-loyalty-api/dto"
	"github.com/gimhanr9/go-loyalty-api/services"
	"github.com/gin-gonic/gin"
)

func RedeemPoints(c *gin.Context) {
	var req dto.RedeemPointsDTO

	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 || req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	req.AccountId = c.GetString("customer_id")

	err := services.RedeemPoints(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rewardTier, err := services.GetDiscountPercentageByClosestRewardTier(req.AccountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	balance, err := services.GetBalance(req.AccountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance, "rewardtier": rewardTier})
}

func EarnPoints(c *gin.Context) {
	var req dto.EarnPointsDTO

	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 || req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	req.AccountId = c.GetString("customer_id")

	err := services.EarnPoints(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rewardTier, err := services.GetDiscountPercentageByClosestRewardTier(req.AccountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	balance, err := services.GetBalance(req.AccountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance, "rewardtier": rewardTier})
}
func GetBalance(c *gin.Context) {
	accountId := c.GetString("customer_id")

	balance, err := services.GetBalance(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func GetHistory(c *gin.Context) {
	accountId := c.GetString("customer_id")
	if accountId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing customer_id in context"})
		return
	}

	cursor := c.Query("cursor") // read from query param

	history, err := services.GetHistory(accountId, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(), // Return error message
		})
		return
	}

	c.JSON(http.StatusOK, history)
}

func GetRewardTiers(c *gin.Context) {
	accountId := c.GetString("customer_id")
	if accountId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing customer_id in context"})
		return
	}

	rewardTier, err := services.GetDiscountPercentageByClosestRewardTier(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	balance, err := services.GetBalance(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance, "rewardtier": rewardTier})
}
