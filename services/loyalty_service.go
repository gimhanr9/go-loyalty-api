package services

import (
	"context"
	"fmt"
	"os"

	square "github.com/square/square-go-sdk"
	client "github.com/square/square-go-sdk/client"
	loyalty "github.com/square/square-go-sdk/loyalty"
	option "github.com/square/square-go-sdk/option"

	"github.com/google/uuid"
)

type LoyaltyService interface {
	EarnPoints(accountID string, description string, LoyaltyEventAccumulatePoints int) error
	//RedeemPoints(accountID string, rewardTierID string) error
	GetBalance(accountID string) (int, error)
	GetHistory(accountID string) ([]square.LoyaltyEvent, error)
}

// EarnPoints adds points to the loyalty account
func EarnPoints(accountID string, description string, amount int) error {

	squareClient := client.NewClient(
		option.WithBaseURL(
			square.Environments.Sandbox,
		),
		option.WithToken(os.Getenv("SQUARE_ACCESS_TOKEN")),
	)

	idempotencyKey := uuid.New().String()

	reqOrder := &square.CreateOrderRequest{
		Order: &square.Order{
			LineItems: []*square.OrderLineItem{
				{
					Name:     &description,
					Quantity: "1",
					BasePriceMoney: &square.Money{
						Amount:   square.Int64(int64(amount)),
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

	orderID := *resOrder.Order.ID

	reqPayment := &square.CreatePaymentRequest{
		SourceID: "cnon:card-nonce-ok",
		AmountMoney: &square.Money{
			Amount:   square.Int64(int64(amount)),
			Currency: square.CurrencyUsd.Ptr(),
		},
		OrderID:        &orderID,
		IdempotencyKey: idempotencyKey,
	}

	_, err = squareClient.Payments.Create(context.TODO(), reqPayment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	reqAccumulate := &loyalty.AccumulateLoyaltyPointsRequest{
		AccountID: accountID,
		AccumulatePoints: &square.LoyaltyEventAccumulatePoints{
			OrderID: square.String(orderID),
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

// // RedeemPoints redeems points for a reward tier
// func (s *loyaltyService) RedeemPoints(accountID string, rewardTierID string, productDescription string, amount int) error {
// 	idempotencyKey := uuid.New().String()

// 	//Create order
// 	reqOrder := &square.CreateOrderRequest{
// 		Order: &square.Order{
// 			LineItems: []*square.OrderLineItem{
// 				&square.OrderLineItem{
// 					Name:     &productDescription,
// 					Quantity: "1",
// 					BasePriceMoney: &square.Money{
// 						Amount: square.Int64(
// 							int64(amount),
// 						),
// 						Currency: square.CurrencyUsd.Ptr(),
// 					},
// 				},
// 			},
// 			LocationID: os.Getenv("LOCATION_ID"),
// 		},
// 		IdempotencyKey: &idempotencyKey,
// 	}

// 	resOrder, errOrder := s.squareClient.Orders.Create(context.TODO(), reqOrder)
// 	if errOrder != nil {
// 		return errors.New("failed to create order")
// 	}

// 	var orderId = *resOrder.Order.ID

// 	reqAccumulate := &loyalty.AccumulateLoyaltyPointsRequest{
// 		AccountID: accountID,
// 		AccumulatePoints: &square.LoyaltyEventAccumulatePoints{
// 			OrderID: square.String(
// 				orderId,
// 			),
// 		},
// 		LocationID:     os.Getenv("LOCATION_ID"),
// 		IdempotencyKey: idempotencyKey,
// 	}
// }

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

// GetHistory retrieves loyalty events (transactions, redemptions, etc.) for the account
func GetHistory(accountID string) ([]square.LoyaltyEvent, error) {

	squareClient := client.NewClient(
		option.WithBaseURL(
			square.Environments.Sandbox,
		),
		option.WithToken(os.Getenv("SQUARE_ACCESS_TOKEN")),
	)

	resp, err := squareClient.Loyalty.SearchEvents(context.TODO(),
		&square.SearchLoyaltyEventsRequest{
			Query: &square.LoyaltyEventQuery{
				Filter: &square.LoyaltyEventFilter{
					LoyaltyAccountFilter: &square.LoyaltyEventLoyaltyAccountFilter{
						LoyaltyAccountID: accountID,
					},
				},
			},
			Limit: square.Int(30),
		})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch loyalty history for account %s: %w", accountID, err)
	}
	if resp == nil || resp.Events == nil {
		return []square.LoyaltyEvent{}, nil
	}

	// Convert []*square.LoyaltyEvent to []square.LoyaltyEvent
	events := make([]square.LoyaltyEvent, len(resp.Events))
	for i, e := range resp.Events {
		if e != nil {
			events[i] = *e
		}
	}

	return events, nil
}
