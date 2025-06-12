package services

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/gimhanr9/go-loyalty-api/dto"
	"github.com/gimhanr9/go-loyalty-api/models"
	"github.com/gimhanr9/go-loyalty-api/repositories"
	"github.com/google/uuid"
	square "github.com/square/square-go-sdk"
	client "github.com/square/square-go-sdk/client"
	loyalty "github.com/square/square-go-sdk/loyalty"
	option "github.com/square/square-go-sdk/option"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dto.RegisterDTO) (*models.User, error)
	Login(req dto.LoginDTO) (*models.User, error)
}

type authService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) AuthService {

	return &authService{
		repo: repo,
	}
}

func (s *authService) Register(req dto.RegisterDTO) (*models.User, error) {
	// Check for existing email or phone
	existing, _ := s.repo.GetByEmailOrPhone(req.Email, req.Phone)
	if existing != nil {
		return nil, errors.New("user with email or phone already exists")
	}

	squareClient := client.NewClient(
		option.WithBaseURL(
			square.Environments.Sandbox,
		),
		option.WithToken(os.Getenv("SQUARE_ACCESS_TOKEN")),
	)

	programRes, programErr := squareClient.Loyalty.Programs.Get(
		context.TODO(),
		&loyalty.GetProgramsRequest{
			ProgramID: "main", //Default program ID
		},
	)

	if programErr != nil {
		return nil, errors.New(programErr.Error())
	}

	programID = *programRes.Program.ID

	// Create Loyalty Account in Square
	idempotencyKey := uuid.New().String()

	reqReg := &loyalty.CreateLoyaltyAccountRequest{
		LoyaltyAccount: &square.LoyaltyAccount{
			Mapping: &square.LoyaltyAccountMapping{
				PhoneNumber: square.String(req.Phone),
			},
			ProgramID: programID,
		},
		IdempotencyKey: idempotencyKey,
	}

	res, err := squareClient.Loyalty.Accounts.Create(context.TODO(), reqReg)
	if err != nil || res.LoyaltyAccount == nil {
		return nil, fmt.Errorf("failed to create loyalty account: %v", err)
	}

	customerId := *res.LoyaltyAccount.ID

	user := &models.User{
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		CustomerID: customerId,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req dto.LoginDTO) (*models.User, error) {
	user, err := s.repo.GetByPhone(req.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}
