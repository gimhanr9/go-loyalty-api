package services

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/square/square-go-sdk/client"
	"github.com/square/square-go-sdk/loyalty"
	"github.com/square/square-go-sdk/option"
)

var (
	clientOnce   sync.Once
	squareClient *client.Client

	programOnce sync.Once
	programID   string
	programErr  error
)

// InitSquareClient initializes the Square client only once and returns it.
func InitSquareClient() *client.Client {
	clientOnce.Do(func() {
		squareClient = client.NewClient(
			option.WithToken(os.Getenv("SQUARE_ACCESS_TOKEN")),
		)
	})
	return squareClient
}

// FetchProgramID retrieves and caches the Square Loyalty Program ID.
func FetchProgramID() (string, error) {
	programOnce.Do(func() {
		client := InitSquareClient()

		resp, err := client.Loyalty.Programs.Get(
			context.TODO(),
			&loyalty.GetProgramsRequest{
				ProgramID: "main", //Default program ID
			},
		)
		if err != nil {
			programErr = fmt.Errorf("error fetching loyalty program: %v", err)
			return
		}

		if resp.Program == nil || resp.Program.ID == nil {
			programErr = fmt.Errorf("no active loyalty program found")
			return
		}

		programID = *resp.Program.ID
	})

	return programID, programErr
}
