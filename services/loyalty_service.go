package services

import (
	"context"
	"fmt"
	"os"
	"time"

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

type Transaction struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Points    int    `json:"points"`
	Timestamp string `json:"timestamp"`
}

// Loyalty History return reponse
type MappedLoyaltyHistoryResponse struct {
	Transactions []Transaction `json:"transactions"`
	Cursor       string        `json:"cursor"`
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

	custoemrRes, err := squareClient.Loyalty.Accounts.Search(
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

	customerId := custoemrRes.LoyaltyAccounts[0].CustomerID

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
			CustomerID: customerId,
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

	paymentRes, err := squareClient.Payments.Create(context.TODO(), reqPayment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	if paymentRes.Payment == nil || paymentRes.Payment.Status == nil || *paymentRes.Payment.Status != "COMPLETED" {
		return fmt.Errorf("payment not completed, status: %v", paymentRes.Payment.Status)
	}

	reqAccumulate := &loyalty.AccumulateLoyaltyPointsRequest{
		AccountID: accountID,
		AccumulatePoints: &square.LoyaltyEventAccumulatePoints{
			OrderID:          square.String(orderID),
			LoyaltyProgramID: &programID,
		},
		LocationID:     os.Getenv("LOCATION_ID"),
		IdempotencyKey: idempotencyKey,
	}

	_, err = squareClient.Loyalty.Accounts.AccumulatePoints(context.TODO(), reqAccumulate)
	if err != nil {
		return fmt.Errorf("failed to accumulate points: %w", err)
	}

	// reqCalculatePoints := &loyalty.CalculateLoyaltyPointsRequest{
	// 	ProgramID: programID,
	// 	TransactionAmountMoney: &square.Money{
	// 		Amount: square.Int64(
	// 			int64(amount),
	// 		),
	// 		Currency: square.CurrencyUsd.Ptr(),
	// 	},
	// 	LoyaltyAccountID: square.String(
	// 		accountID,
	// 	),
	// }

	// resCalc, err := squareClient.Loyalty.Programs.Calculate(context.TODO(), reqCalculatePoints)

	// if err != nil {
	// 	return fmt.Errorf("failed to calculate points: %w", err)
	// }

	// pointsEarned := *resCalc.Points

	// if pointsEarned > 0 {
	// 	reqAccumulate := &loyalty.AccumulateLoyaltyPointsRequest{
	// 		AccountID: accountID,
	// 		AccumulatePoints: &square.LoyaltyEventAccumulatePoints{
	// 			OrderID: square.String(orderID),
	// 		},
	// 		LocationID:     os.Getenv("LOCATION_ID"),
	// 		IdempotencyKey: idempotencyKey,
	// 	}

	// 	_, err = squareClient.Loyalty.Accounts.AccumulatePoints(context.TODO(), reqAccumulate)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to accumulate points: %w", err)
	// 	}
	// }

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

func formatTimestamp(raw string) string {
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}
	return t.Format("2 Jan 2006 15:04")
}

// GetHistory retrieves loyalty events (transactions, redemptions, etc.) for the account
func GetHistory(accountID string, cursor string) (*MappedLoyaltyHistoryResponse, error) {
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

	transactions := make([]Transaction, 0)
	if resp.Events != nil {
		for _, e := range resp.Events {
			if e == nil {
				continue
			}

			points := 0
			if e.AccumulatePoints != nil && e.AccumulatePoints.Points != nil {
				points = int(*e.AccumulatePoints.Points)
			}

			transactions = append(transactions, Transaction{
				ID:        e.ID,
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

	return &MappedLoyaltyHistoryResponse{
		Transactions: transactions,
		Cursor:       newCursor,
	}, nil
}
