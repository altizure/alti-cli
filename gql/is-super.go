package gql

import (
	"context"

	"github.com/machinebox/graphql"
)

// IsSuper checks if the provided creds is a superuser.
func IsSuper(endpoint, key, token string) bool {
	client := graphql.NewClient(endpoint + "/super")

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

type isSuperRes struct {
	Hello string
}
