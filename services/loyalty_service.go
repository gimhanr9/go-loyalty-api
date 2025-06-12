package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	square "github.com/square/square-go-sdk"
	client "github.com/square/square-go-sdk/client"
	loyalty "github.com/square/square-go-sdk/loyalty"
	option "github.com/square/square-go-sdk/option"

	"github.com/google/uuid"

	"github.com/gimhanr9/go-loyalty-api/dto"
)

type LoyaltyService interface {
	EarnPoints(req dto.EarnPointsDTO) error
	RedeemPoints(req dto.RedeemPointsDTO) error
	GetBalance(accountID string) (int, error)
	GetHistory(accountID string) ([]square.LoyaltyEvent, error)
	GetDiscountPercentageByClosestRewardTier(accountID string) (*dto.RewardTierDTO, error)
}

// EarnPoints adds points to the loyalty account
func EarnPoints(req dto.EarnPointsDTO) error {

	squareClient := client.NewClient(
		option.WithBaseURL(
			square.Environments.Sandbox,
		),
		option.WithToken(os.Getenv("SQUARE_ACCESS_TOKEN")),
	)

	idempotencyKey := uuid.New().String()

	//Get program Id
	programRes, err := squareClient.Loyalty.Programs.Get(
		context.TODO(),
		&loyalty.GetProgramsRequest{
			ProgramID: "main", //Default program ID
		},
	)

	if err != nil {
		return fmt.Errorf("failed to retrieve program: %w", err)
	}

	programID = *programRes.Program.ID

	customerRes, err := squareClient.Loyalty.Accounts.Search(
		context.TODO(),
		&loyalty.SearchLoyaltyAccountsRequest{
			Limit: square.Int(
				1,
			),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	customerId := customerRes.LoyaltyAccounts[0].CustomerID

	reqOrder := &square.CreateOrderRequest{
		Order: &square.Order{
			LineItems: []*square.OrderLineItem{
				{
					Name:     square.String(req.Description),
					Quantity: "1",
					BasePriceMoney: &square.Money{
						Amount:   square.Int64(int64(req.Amount)),
						Currency: square.CurrencyUsd.Ptr(),
					},
				},
			},
			CustomerID: customerId,
			LocationID: os.Getenv("LOCATION_ID"),
		},

		IdempotencyKey: &idempotencyKey,
	}

	resOrder, err := squareClient.Orders.Create(context.TODO(), reqOrder)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	orderId := *resOrder.Order.ID

	reqPayment := &square.CreatePaymentRequest{
		SourceID: "cnon:card-nonce-ok",
		AmountMoney: &square.Money{
			Amount:   square.Int64(int64(req.Amount)),
			Currency: square.CurrencyUsd.Ptr(),
		},
		OrderID:        square.String(orderId),
		IdempotencyKey: idempotencyKey,
	}

	paymentRes, err := squareClient.Payments.Create(context.TODO(), reqPayment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	if paymentRes.Payment == nil || paymentRes.Payment.Status == nil || *paymentRes.Payment.Status != "COMPLETED" {
		return fmt.Errorf("payment not completed, status: %v", paymentRes.Payment.Status)
	}

	reqAccumulate := &loyalty.AccumulateLoyaltyPointsRequest{
		AccountID: req.AccountId,
		AccumulatePoints: &square.LoyaltyEventAccumulatePoints{
			OrderID: square.String(orderId),
		},
		LocationID:     os.Getenv("LOCATION_ID"),
		IdempotencyKey: idempotencyKey,
	}

	_, err = squareClient.Loyalty.Accounts.AccumulatePoints(context.TODO(), reqAccumulate)
	if err != nil {
		return fmt.Errorf("failed to accumulate points: %w", err)
	}

	return nil
}

// RedeemPoints redeems points for a reward tier
func RedeemPoints(req dto.RedeemPointsDTO) error {

	squareClient := client.NewClient(
		option.WithBaseURL(
			square.Environments.Sandbox,
		),
		option.WithToken(os.Getenv("SQUARE_ACCESS_TOKEN")),
	)

	idempotencyKey := uuid.New().String()

	//Create order
	reqOrder := &square.CreateOrderRequest{
		Order: &square.Order{
			LineItems: []*square.OrderLineItem{
				&square.OrderLineItem{
					Name:     square.String(req.Description),
					Quantity: "1",
					BasePriceMoney: &square.Money{
						Amount: square.Int64(
							int64(req.Amount),
						),
						Currency: square.CurrencyUsd.Ptr(),
					},
				},
			},
			LocationID: os.Getenv("LOCATION_ID"),
		},
		IdempotencyKey: &idempotencyKey,
	}

	resOrder, err := squareClient.Orders.Create(context.TODO(), reqOrder)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	var orderId = *resOrder.Order.ID

	reqReward := &loyalty.CreateLoyaltyRewardRequest{
		Reward: &square.LoyaltyReward{
			LoyaltyAccountID: req.AccountId,
			RewardTierID:     req.RewardTierId,
			OrderID: square.String(
				orderId,
			),
		},
		IdempotencyKey: idempotencyKey,
	}

	_, err = squareClient.Loyalty.Rewards.Create(context.TODO(), reqReward)
	if err != nil {
		return errors.New("failed to create reward")
	}

	discountOrderRes, err := squareClient.Orders.Get(
		context.TODO(),
		&square.GetOrdersRequest{
			OrderID: orderId,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to get order details: %w", err)
	}

	reqPayment := &square.CreatePaymentRequest{
		SourceID: "cnon:card-nonce-ok",
		AmountMoney: &square.Money{
			Amount:   discountOrderRes.Order.TotalMoney.Amount,
			Currency: square.CurrencyUsd.Ptr(),
		},
		OrderID:        square.String(orderId),
		IdempotencyKey: idempotencyKey,
	}

	paymentRes, err := squareClient.Payments.Create(context.TODO(), reqPayment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	if paymentRes.Payment == nil || paymentRes.Payment.Status == nil || *paymentRes.Payment.Status != "COMPLETED" {
		return fmt.Errorf("payment not completed, status: %v", paymentRes.Payment.Status)
	}

	reqAccumulate := &loyalty.AccumulateLoyaltyPointsRequest{
		AccountID: req.AccountId,
		AccumulatePoints: &square.LoyaltyEventAccumulatePoints{
			OrderID:          square.String(orderId),
			LoyaltyProgramID: &programID,
		},
		LocationID:     os.Getenv("LOCATION_ID"),
		IdempotencyKey: idempotencyKey,
	}

	_, err = squareClient.Loyalty.Accounts.AccumulatePoints(context.TODO(), reqAccumulate)
	if err != nil {
		return fmt.Errorf("failed to accumulate points: %w", err)
	}

	return nil
}

// GetBalance fetches the points balance of the loyalty account
func GetBalance(accountID string) (int, error) {

	squareClient := client.NewClient(
		option.WithBaseURL(
			square.Environments.Sandbox,
		),
		option.WithToken(os.Getenv("SQUARE_ACCESS_TOKEN")),
	)

	resp, err := squareClient.Loyalty.Accounts.Get(context.TODO(),
		&loyalty.GetAccountsRequest{
			AccountID: accountID,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get account %s balance: %w", accountID, err)
	}

	if resp == nil || resp.LoyaltyAccount == nil || resp.LoyaltyAccount.Balance == nil {
		return 0, fmt.Errorf("no balance information found for account %s", accountID)
	}

	return *resp.LoyaltyAccount.Balance, nil
}

func formatTimestamp(raw string) string {
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}
	return t.Format("2 Jan 2006 15:04")
}

// GetHistory retrieves loyalty events (transactions, redemptions, etc.) for the account
func GetHistory(accountID string, cursor string) (*dto.MappedLoyaltyHistoryResponseDTO, error) {
	squareClient := client.NewClient(
		option.WithBaseURL(square.Environments.Sandbox),
		option.WithToken(os.Getenv("SQUARE_ACCESS_TOKEN")),
	)

	req := &square.SearchLoyaltyEventsRequest{
		Query: &square.LoyaltyEventQuery{
			Filter: &square.LoyaltyEventFilter{
				LoyaltyAccountFilter: &square.LoyaltyEventLoyaltyAccountFilter{
					LoyaltyAccountID: accountID,
				},
			},
		},
		Limit: square.Int(10),
	}

	if cursor != "" {
		req.Cursor = square.String(cursor)
	}

	resp, err := squareClient.Loyalty.SearchEvents(context.TODO(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch loyalty history for account %s: %w", accountID, err)
	}

	transactions := make([]dto.TransactionDTO, 0)
	if resp.Events != nil {
		for _, e := range resp.Events {
			if e == nil {
				continue
			}

			points := 0
			if e.AccumulatePoints != nil && e.AccumulatePoints.Points != nil {
				points = int(*e.AccumulatePoints.Points)
			}

			transactions = append(transactions, dto.TransactionDTO{
				Id:        e.ID,
				Type:      string(e.Type),
				Points:    points,
				Timestamp: formatTimestamp(e.CreatedAt),
			})
		}
	}

	newCursor := ""
	if c := resp.GetCursor(); c != nil {
		newCursor = *c
	}

	return &dto.MappedLoyaltyHistoryResponseDTO{
		Transactions: transactions,
		Cursor:       newCursor,
	}, nil
}

func MapClosestRewardTier(program *square.LoyaltyProgram, userBalance int) *dto.RewardTierDTO {
	var closestTier *square.LoyaltyProgramRewardTier
	closestPoints := -1

	for _, tier := range program.RewardTiers {
		if tier != nil {
			points := tier.Points

			// Choose the tier with highest required points ≤ userBalance
			if points <= userBalance && points > closestPoints {
				closestPoints = points
				closestTier = tier
			}
		}
	}

	if closestTier == nil || closestTier.ID == nil {
		return nil // No valid tier found
	}

	percentageStr := closestTier.Definition.PercentageDiscount
	var discountPercentage float32

	if percentageStr != nil {
		if f, err := strconv.ParseFloat(*percentageStr, 32); err == nil {
			discountPercentage = float32(f)
		}
	}

	return &dto.RewardTierDTO{
		RewardTierId:       *closestTier.ID,
		DiscountPercentage: discountPercentage,
	}
}

// GetHistory retrieves loyalty events (transactions, redemptions, etc.) for the account
func GetDiscountPercentageByClosestRewardTier(accountID string) (*dto.RewardTierDTO, error) {
	squareClient := client.NewClient(
		option.WithBaseURL(square.Environments.Sandbox),
		option.WithToken(os.Getenv("SQUARE_ACCESS_TOKEN")),
	)

	// Get loyalty account
	resp, err := squareClient.Loyalty.Accounts.Get(context.TODO(),
		&loyalty.GetAccountsRequest{AccountID: accountID})
	if err != nil || resp.LoyaltyAccount == nil {
		return &dto.RewardTierDTO{RewardTierId: "", DiscountPercentage: 0}, fmt.Errorf("failed to get account %s balance: %w", accountID, err)
	}

	// Safe balance
	balance := 0
	if resp.LoyaltyAccount.Balance != nil {
		balance = *resp.LoyaltyAccount.Balance
	}

	// Get loyalty program
	programRes, err := squareClient.Loyalty.Programs.Get(context.TODO(),
		&loyalty.GetProgramsRequest{ProgramID: "main"})
	if err != nil || programRes.Program == nil {
		return &dto.RewardTierDTO{RewardTierId: "", DiscountPercentage: 0}, fmt.Errorf("failed to retrieve loyalty program: %w", err)
	}

	// Get closest tier
	rewardTier := MapClosestRewardTier(programRes.Program, balance)
	if rewardTier == nil {
		return &dto.RewardTierDTO{RewardTierId: "", DiscountPercentage: 0}, nil
	}

	return rewardTier, nil
}
