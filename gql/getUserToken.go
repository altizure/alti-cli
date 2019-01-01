package gql

import (
	"context"
	"log"

	"github.com/machinebox/graphql"
)

// GetUserToken gets the self-issued user token.
func GetUserToken(endpoint, email, password string, fresh bool) string {
	// create a client (safe to share across requests)
	client := graphql.NewClient(endpoint)

	// make a request
	req := graphql.NewRequest(`
		mutation ($email: String!, $password: String!) {
		  getUserToken(email: $email, password: $password)
		}
	`)

	// set any variables
	req.Var("email", email)
	req.Var("password", password)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res response
	if err := client.Run(ctx, req, &res); err != nil {
		log.Fatal(err)
	}
	return res.GetUserToken
}

type response struct {
	GetUserToken string
}
