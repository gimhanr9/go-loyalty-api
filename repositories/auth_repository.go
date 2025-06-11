package repositories

import (
	"github.com/gimhanr9/go-loyalty-api/database"
	"github.com/gimhanr9/go-loyalty-api/models"
)

type AuthRepository interface {
	GetByEmailOrPhone(email, phone string) (*models.User, error)
	GetByPhone(phone string) (*models.User, error)
	Create(user *models.User) error
}

type authRepository struct{}

func NewAuthRepository() AuthRepository {
	return &authRepository{}
}

func (r *authRepository) GetByEmailOrPhone(email, phone string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ? OR phone = ?", email, phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetByPhone(phone string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}
