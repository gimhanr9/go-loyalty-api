package services

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/gimhanr9/go-loyalty-api/models"
	"github.com/gimhanr9/go-loyalty-api/repositories"
	square "github.com/square/square-go-sdk"
	"github.com/square/square-go-sdk/loyalty"
	"github.com/square/square-go-sdk/option"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(name, email, phone string) (*models.User, error)
	Login(phone string) (*models.User, error)
}

type authService struct {
	repo         repositories.AuthRepository
	squareClient *square.Client
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	client := square.NewClient(option.WithToken(os.Getenv("SQUARE_ACCESS_TOKEN")))

	return &authService{
		repo:         repo,
		squareClient: client,
	}
}

func (s *authService) Register(name, email, phone string) (*models.User, error) {
	// Check for existing email or phone
	existing, _ := s.repo.GetByEmailOrPhone(email, phone)
	if existing != nil {
		return nil, errors.New("user with email or phone already exists")
	}

	programID, err := FetchProgramID()
	if err != nil {
		return nil, errors.New("failed to fetch loyalty program ID")
	}

	// Create in Square
	idempotencyKey := uuid.New().String()

	req := &loyalty.CreateLoyaltyAccountRequest{
		LoyaltyAccount: &square.LoyaltyAccount{
			Mapping: &square.LoyaltyAccountMapping{
				PhoneNumber: &phone,
			},
			ProgramID: programID,
		},
		IdempotencyKey: idempotencyKey,
	}

	res, err := s.squareClient.Loyalty.Accounts.Create(req)
	if err != nil || res.LoyaltyAccount == nil {
		return nil, fmt.Errorf("failed to create loyalty account: %v", err)
	}

	customerID := *res.LoyaltyAccount.ID

	user := &models.User{
		Name:       name,
		Email:      email,
		Phone:      phone,
		CustomerID: customerID,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(phone string) (*models.User, error) {
	user, err := s.repo.GetByPhone(phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}
