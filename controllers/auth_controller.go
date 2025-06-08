package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gimhanr9/go-loyalty-api/services"
	"net/http"
)

func Login(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.CustomerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CustomerID required"})
		return
	}

	token, err := services.AuthService{}.Login(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
