package gql

import (
	"context"
	"errors"

	"github.com/machinebox/graphql"
)

// RequestLoginCode requests an one-time login code via sms.
func RequestLoginCode(endpoint, appKey, phone string) error {
	client := graphql.NewClient(endpoint + "/graphql")

	req := graphql.NewRequest(`
		mutation ($phone: String!) {
			requestLoginCode(phone: $phone) {
				error {
					code
					message
				}
				result
			}
		}
	`)
	req.Header.Set("key", appKey)

	req.Var("phone", phone)

	ctx := context.Background()
	var res reqLoginCodeRes
	if err := client.Run(ctx, req, &res); err != nil {
		return err
	}
	if res.RequestLoginCode.Result != "Success" {
		return errors.New(res.RequestLoginCode.Error.Message)
	}
	return nil
}

type reqLoginCodeRes struct {
	RequestLoginCode struct {
		Error struct {
			Code    int
			Message string
		}
		Result string
	}
}
