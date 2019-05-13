package gql

import (
	"context"
	"log"

	"github.com/machinebox/graphql"
)

// GetUserTokenByCode gets the self-issued by phone and one-time login code.
func GetUserTokenByCode(endpoint, appKey, phone, code string) string {
	client := graphql.NewClient(endpoint + "/graphql")

	req := graphql.NewRequest(`
		mutation ($phone: String!, $code: String!) {
			getUserTokenByLoginCode(phone: $phone, code: $code, fresh: false)
		}
	`)
	req.Header.Set("key", appKey)

	req.Var("phone", phone)
	req.Var("code", code)

	ctx := context.Background()

	// run it and capture the response
	var res getUserTokenByCodeRes
	if err := client.Run(ctx, req, &res); err != nil {
		log.Fatal(err)
	}
	return res.GetUserTokenByLoginCode
}

type getUserTokenByCodeRes struct {
	GetUserTokenByLoginCode string
}
