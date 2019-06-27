package gql

import (
	"context"

	"github.com/machinebox/graphql"
)

// GetUserTokenByEmail gets the self-issued user token.
func GetUserTokenByEmail(endpoint, appKey, email, password string, fresh bool) (string, error) {
	// create a client (safe to share across requests)
	client := graphql.NewClient(endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($email: String!, $password: String!) {
			getUserToken(email: $email, password: $password, fresh: false)
		}
	`)
	req.Header.Set("key", appKey)

	// set any variables
	req.Var("email", email)
	req.Var("password", password)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res getUserTokenRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	return res.GetUserToken, nil
}

type getUserTokenRes struct {
	GetUserToken string
}
