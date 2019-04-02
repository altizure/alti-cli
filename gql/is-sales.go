package gql

import (
	"context"

	"github.com/machinebox/graphql"
)

// IsSales checks if the provided creds is a sales.
func IsSales(endpoint, key, token string) bool {
	client := graphql.NewClient(endpoint + "/sales")

	req := graphql.NewRequest(`
		{
		  hello
		}
	`)
	req.Header.Set("key", key)
	req.Header.Set("altitoken", token)

	ctx := context.Background()
	var res isSalesRes
	if err := client.Run(ctx, req, &res); err != nil {
		return false
	}
	if res.Hello == "" {
		return false
	}
	return true
}

type isSalesRes struct {
	Hello string
}
