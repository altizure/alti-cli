package gql

import (
	"context"
	"errors"

	"github.com/jackytck/alti-cli/config"
	altiErrors "github.com/jackytck/alti-cli/errors"
	"github.com/machinebox/graphql"
)

// TransferCoins transfers coins from my account to other user,
// with a custom message.
func TransferCoins(coins float64, email, message string) (string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(`
		mutation ($amount: Float!, $email: String!, $message: String){
			transferCoins(amount: $amount, options: {email: $email, message: $message}) {
				error {
					message
				}
				result
			}
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// set variables
	req.Var("amount", coins)
	req.Var("email", email)
	req.Var("message", message)

	ctx := context.Background()

	var res transCoinsRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	errMsg := res.TransferCoins.Error.Message
	result := res.TransferCoins.Result
	if errMsg != "" {
		return result, errors.New(errMsg)
	}
	if result == "Fail" {
		return result, altiErrors.ErrTransferCoins
	}

	return result, nil
}

type transCoinsRes struct {
	TransferCoins struct {
		Error struct {
			Message string
		}
		Result string
	}
}
