package repositories

import (
	"github.com/gimhanr9/go-loyalty-api/database"
	"github.com/gimhanr9/go-loyalty-api/models"
)

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	return &user, result.Error
}

func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}
