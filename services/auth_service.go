package services

import (
	"errors"
	"github.com/gimhanr9/go-loyalty-api/repositories"
	"github.com/gimhanr9/go-loyalty-api/utils"
)

type AuthService struct{}

func (s AuthService) Login(customerID string) (string, error) {
	authRepo := repositories.AuthRepository{}

	exists := authRepo.CheckCustomerExists(customerID)
	if !exists {
		return "", errors.New("customer not found")
	}

	token, err := utils.GenerateJWT(customerID)
	if err != nil {
		return "", err
	}

	return token, nil
}
