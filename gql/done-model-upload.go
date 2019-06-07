package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/machinebox/graphql"
)

// DoneModelUpload signals the completion of (multipart) model upload.
// Return the state of the project.
func DoneModelUpload(pid string) (string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(`
		mutation ($pid: ID!) {
		  doneModelUpload(id: $pid) {
		    id
		    importedState
		  }
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// set variables
	req.Var("pid", pid)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res doneModelUploadRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	id := res.DoneModelUpload.ID
	state := res.DoneModelUpload.ImportedState
	if id == "" {
		return state, errors.ErrModelMutateState
	}

	return state, nil
}

type doneModelUploadRes struct {
	DoneModelUpload struct {
		ID            string
		ImportedState string
	}
}
