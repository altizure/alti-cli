package gql

import (
	"context"
	"errors"

	"github.com/jackytck/alti-cli/config"
	altiErrors "github.com/jackytck/alti-cli/errors"
	"github.com/machinebox/graphql"
)

// TransferProject transfers project from my account to other user,
// with a custom message.
func TransferProject(pid, email, message string) (string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(`
		mutation ($id: ID!, $email: String!, $message: String){
			transferProject(id: $id, options: {email: $email, message: $message}) {
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
	req.Var("id", pid)
	req.Var("email", email)
	req.Var("message", message)

	ctx := context.Background()

	var res transProjRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	errMsg := res.TransferProject.Error.Message
	result := res.TransferProject.Result
	if errMsg != "" {
		return result, errors.New(errMsg)
	}
	if result == "Fail" {
		return result, altiErrors.ErrTransferProject
	}

	return result, nil
}

type transProjRes struct {
	TransferProject struct {
		Error struct {
			Message string
		}
		Result string
	}
}
