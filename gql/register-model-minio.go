package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RegisterModelMinio registers a Minio model.
// And get back the registered model and the signed url to Minio.
func RegisterModelMinio(pid, bucket, filename string) (*types.Model, string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($id: ID!, $bucket: BucketMinioModel!, $filename: String!) {
			uploadModelMinio(id: $id, bucket: $bucket, filename: $filename) {
				url
				file {
					id
					state
					name
					filename
				}
			}
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// set variables
	req.Var("id", pid)
	req.Var("bucket", bucket)
	req.Var("filename", filename)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res regModelMinioRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, "", err
	}
	mid := res.UploadModelMinio.File.ID
	url := res.UploadModelMinio.URL
	if mid == "" || url == "" {
		return nil, "", errors.ErrModelReg
	}

	return &res.UploadModelMinio.File, url, nil
}

type regModelMinioRes struct {
	UploadModelMinio struct {
		URL  string
		File types.Model
	}
}
