package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RegisterImageURL registers an to be uploaded image by url.
func RegisterImageURL(pid, url, filename, checksum string) (*types.Image, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($pid: ID!, $url: String!, $filename: String, $checksum: String) {
		  uploadImageURL(pid: $pid, url: $url, filename: $filename, checksum: $checksum) {
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
	var res regImgURLRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	iid := res.uploadImageURL.ID
	if iid == "" {
		return nil, errors.ErrImgReg
	}

	img := types.Image{
		ID:       res.uploadImageURL.ID,
		State:    res.uploadImageURL.State,
		Name:     res.uploadImageURL.Name,
		Filename: res.uploadImageURL.Filename,
	}

	return &img, nil
}

type regImgURLRes struct {
	uploadImageURL struct {
		ID       string
		State    string
		Name     string
		Filename string
	}
}
