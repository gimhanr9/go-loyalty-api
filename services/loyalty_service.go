package services

import (
	"context"
	"errors"
	"os"

	square "github.com/square/square-go-sdk"
	client "github.com/square/square-go-sdk/client"
	loyalty "github.com/square/square-go-sdk/loyalty"
	option "github.com/square/square-go-sdk/option"

	"github.com/google/uuid"
)

type LoyaltyService interface {
	EarnPoints(accountID string, points int) error
	RedeemPoints(accountID string, rewardTierID string) error
	GetBalance(accountID string) (int, error)
	GetHistory(accountID string) ([]loyalty.LoyaltyEvent, error)
}

type loyaltyService struct {
	squareClient *client.Client
}

func NewLoyaltyService() LoyaltyService {
	squareClient := client.NewClient(
		option.WithToken(os.Getenv("SQUARE_ACCESS_TOKEN")),
	)

	return &loyaltyService{
		squareClient: squareClient,
	}
}

func generateUUID() string {
	return uuid.New().String()
}

// EarnPoints adds points to the loyalty account
func (s *loyaltyService) EarnPoints(accountID string, productDescription string, amount int) error {

	idempotencyKey := uuid.New().String()

	//Create order
	reqOrder := &square.CreateOrderRequest{
		Order: &square.Order{
			LineItems: []*square.OrderLineItem{
				&square.OrderLineItem{
					Name:     &productDescription,
					Quantity: "1",
					BasePriceMoney: &square.Money{
						Amount: square.Int64(
							int64(amount),
						),
						Currency: square.CurrencyUsd.Ptr(),
					},
				},
			},
			LocationID: os.Getenv("LOCATION_ID"),
		},
		IdempotencyKey: &idempotencyKey,
	}

	resOrder, errOrder := s.squareClient.Orders.Create(context.TODO(), reqOrder)
	if errOrder != nil {
		return errors.New("failed to create order")
	}

	var orderId = *resOrder.Order.ID

	//Create payment
	reqPayment := &square.CreatePaymentRequest{
		SourceID: "cnon:card-nonce-ok",
		AmountMoney: &square.Money{
			Amount: square.Int64(
				int64(amount),
			),
			Currency: square.CurrencyUsd.Ptr(),
		},
		OrderID:        &orderId,
		IdempotencyKey: idempotencyKey,
	}

	_, errPayment := s.squareClient.Payments.Create(context.TODO(), reqPayment)
	if errPayment != nil {
		return errors.New("failed to create payment")
	}

	// Add points to loyalty account
	reqAccumulate := &loyalty.AccumulateLoyaltyPointsRequest{
		AccountID: accountID,
		AccumulatePoints: &square.LoyaltyEventAccumulatePoints{
			OrderID: square.String(
				orderId,
			),
		},
		LocationID:     os.Getenv("LOCATION_ID"),
		IdempotencyKey: idempotencyKey,
	}

	_, errAccumulate := s.squareClient.Loyalty.Accounts.AccumulatePoints(context.TODO(), reqAccumulate)
	if errAccumulate != nil {
		return errors.New("failed to accumulate points")
	}
}

// RedeemPoints redeems points for a reward tier
func (s *loyaltyService) RedeemPoints(accountID string, rewardTierID string, productDescription string, amount int) error {
	idempotencyKey := uuid.New().String()

	//Create order
	reqOrder := &square.CreateOrderRequest{
		Order: &square.Order{
			LineItems: []*square.OrderLineItem{
				&square.OrderLineItem{
					Name:     &productDescription,
					Quantity: "1",
					BasePriceMoney: &square.Money{
						Amount: square.Int64(
							int64(amount),
						),
						Currency: square.CurrencyUsd.Ptr(),
					},
				},
			},
			LocationID: os.Getenv("LOCATION_ID"),
		},
		IdempotencyKey: &idempotencyKey,
	}

	resOrder, errOrder := s.squareClient.Orders.Create(context.TODO(), reqOrder)
	if errOrder != nil {
		return errors.New("failed to create order")
	}
	
	var orderId = *resOrder.Order.ID

	reqAccumulate := &loyalty.AccumulateLoyaltyPointsRequest{
		AccountID: accountID,
		AccumulatePoints: &square.LoyaltyEventAccumulatePoints{
			OrderID: square.String(
				orderId,
			),
		},
		LocationID:     os.Getenv("LOCATION_ID"),
		IdempotencyKey: idempotencyKey,
	}
}

// GetBalance fetches the points balance of the loyalty account
func (s *loyaltyService) GetBalance(accountID string) (int, error) {
	resp, err := s.squareClient.Loyalty.Accounts.Get(context.TODO(),
		&loyalty.GetAccountsRequest{
            AccountID: accountID,
        },
	)
	if err != nil {
		return 0, err
	}

	return *resp.LoyaltyAccount.Balance, nil
}

// GetHistory retrieves loyalty events (transactions, redemptions, etc.) for the account
func (s *loyaltyService) GetHistory(accountID string) ([]square.LoyaltyEvent, error) {
	resp, err := s.squareClient.Loyalty.SearchEvents(context.TODO(),
		&square.SearchLoyaltyEventsRequest{
			Query: &square.LoyaltyEventQuery{

				Filter: &square.LoyaltyEventFilter{
					LoyaltyAccountFilter: &square.LoyaltyEventLoyaltyAccountFilter{
						LoyaltyAccountID: accountID,
					},
				},
			},
			Limit: square.Int(
				30,
			),
		})
	if err != nil {
		return nil, err
	}
	return resp.Events, nil
}
