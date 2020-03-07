package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/machinebox/graphql"
)

// SetProfileFace set the profile image with the given image string.
// Return the result of operation.
func SetProfileFace(imgStr string) (string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(`
		mutation ($imgStr: String!) {
			setProfileFace(imgStr: $imgStr)
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// set variables
	req.Var("imgStr", imgStr)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res setProfileFaceRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	return res.SetProfileFace, nil
}

type setProfileFaceRes struct {
	SetProfileFace string
}
