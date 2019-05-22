package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RegisterModelURL registers a to be uploaded model by url.
func RegisterModelURL(pid, url, filename, checksum string) (*types.ImportedModel, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($pid: ID!, $url: String!, $filename: String, $checksum: String) {
		  uploadModelURL(pid: $pid, url: $url, filename: $filename, checksum: $checksum) {
		    id
				state
				name
				filename
		  }
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// set variables
	req.Var("pid", pid)
	req.Var("url", url)
	req.Var("filename", filename)
	req.Var("checksum", checksum)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res regModelURLRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	iid := res.UploadModelURL.ID
	if iid == "" {
		return nil, errors.ErrModelReg
	}

	model := types.ImportedModel{
		ID:       res.UploadModelURL.ID,
		State:    res.UploadModelURL.State,
		Name:     res.UploadModelURL.Name,
		Filename: res.UploadModelURL.Filename,
	}

	return &model, nil
}

type regModelURLRes struct {
	UploadModelURL struct {
		ID       string
		State    string
		Name     string
		Filename string
	}
}
