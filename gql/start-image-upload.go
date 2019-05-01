package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/machinebox/graphql"
)

// StartImageUpload signals the start of image uploading.
// Return the new image state with error.
func StartImageUpload(iid string) (string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(`
		mutation ($iid: ID!) {
			startImageUpload(id: $iid) {
				id
				state
			}
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// set variables
	req.Var("iid", iid)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res startImgUploadRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	id := res.StartImageUpload.ID
	state := res.StartImageUpload.State
	if id == "" {
		return state, errors.ErrImgMutateState
	}

	return state, nil
}

type startImgUploadRes struct {
	StartImageUpload struct {
		ID    string
		State string
	}
}
