package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RegisterMetaURL registers a to be uploaded meta file by url.
func RegisterMetaURL(pid, url, filename, checksum string) (*types.MetaFile, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($pid: ID!, $url: String!, $filename: String, $checksum: String) {
			uploadMetaURL(pid: $pid, url: $url, filename: $filename, checksum: $checksum) {
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
	var res regMetaURLRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	iid := res.UploadMetaURL.ID
	if iid == "" {
		return nil, errors.ErrMetaReg
	}

	return &res.UploadMetaURL, nil
}

type regMetaURLRes struct {
	UploadMetaURL types.MetaFile
}
