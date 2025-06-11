package services

import (
	"context"
	"fmt"
	"os"
	"sync"

	square "github.com/square/square-go-sdk"
	"github.com/square/square-go-sdk/option"
)

var (
	clientOnce  sync.Once
	sqClient    *square.Client

	programID    string
	programOnce  sync.Once
	programErr   error
)

// Init initializes Square client once
func Init() *square.Client {
	clientOnce.Do(func() {
		sqClient = square.NewClient(
			option.WithToken(os.Getenv("SQUARE_ACCESS_TOKEN")),
		)
	})
	return sqClient
}

func FetchProgramID() (string, error) {
	programOnce.Do(func() {
		client := Init()

		resp, err := client.Loyalty.Programs.Get(
			context.TODO(),
			&square.GetLoyaltyProgramRequest{
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
