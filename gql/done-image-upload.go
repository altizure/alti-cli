package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/machinebox/graphql"
)

// DoneImageUpload signals the end of image uploading.
// Return the new image state with error.
func DoneImageUpload(iid string) (string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(`
		mutation ($iid: ID!) {
			doneImageUpload(id: $iid) {
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
	var res doneImgUploadRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	id := res.DoneImageUpload.ID
	state := res.DoneImageUpload.State
	if id == "" {
		return state, errors.ErrImgMutateState
	}

	return state, nil
}

type doneImgUploadRes struct {
	DoneImageUpload struct {
		ID    string
		State string
	}
}
