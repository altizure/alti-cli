package gql

import (
	"context"

	"github.com/machinebox/graphql"
)

// GetUserTokenBySuper gets the self-issued user token by superuser.
func GetUserTokenBySuper(endpoint, appKey, email string) (string, error) {
	client := graphql.NewClient(endpoint + "/super")

	req := graphql.NewRequest(`
		mutation ($email: String!) {
			getUserToken(email: $email)
		}
	`)
	req.Header.Set("key", appKey)
	req.Var("email", email)

	ctx := context.Background()

	var res getUserTokenRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	return res.GetUserToken, nil
}
