package supergql

import (
	"context"
)

// GetUserToken gets the self-issued user token.
func GetUserToken(email string) (string, error) {
	gql := `
		mutation ($email: String!) {
			getUserToken(email: $email)
		}
	`
	req, client := SuperRequest(gql)
	req.Var("email", email)

	ctx := context.Background()

	var res getUserTokenRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	return res.GetUserToken, nil
}

type getUserTokenRes struct {
	GetUserToken string
}
