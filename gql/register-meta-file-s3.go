package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RegisterMetaFileS3 registers a S3 meta file.
// And get back the registered meta file and the signed url to S3.
func RegisterMetaFileS3(pid, bucket, filename string) (*types.MetaFile, string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($pid: ID!, $bucket: BucketS3!, $filename: String!) {
			uploadMetaFileS3(pid: $pid, bucket: $bucket, filename: $filename) {
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
	req.Var("pid", pid)
	req.Var("bucket", bucket)
	req.Var("filename", filename)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res regMetaS3Res
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, "", err
	}
	mid := res.UploadMetaFileS3.File.ID
	url := res.UploadMetaFileS3.URL
	if mid == "" || url == "" {
		return nil, "", errors.ErrModelReg
	}

	return &res.UploadMetaFileS3.File, url, nil
}

type regMetaS3Res struct {
	UploadMetaFileS3 struct {
		URL  string
		File types.MetaFile
	}
}
