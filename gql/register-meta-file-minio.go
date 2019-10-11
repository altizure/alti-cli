package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RegisterMetaFileMinio registers a Minio meta file.
// And get back the registered meta file and the signed url to Minio.
func RegisterMetaFileMinio(pid, bucket, filename string) (*types.MetaFile, string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($pid: ID!, $bucket: BucketMinioMeta!, $filename: String!) {
			uploadMetaFileMinio(pid: $pid, bucket: $bucket, filename: $filename) {
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
	var res regMetaMinioRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, "", err
	}
	mid := res.UploadMetaFileMinio.File.ID
	url := res.UploadMetaFileMinio.URL
	if mid == "" || url == "" {
		return nil, "", errors.ErrMetaReg
	}

	return &res.UploadMetaFileMinio.File, url, nil
}

type regMetaMinioRes struct {
	UploadMetaFileMinio struct {
		URL  string
		File types.MetaFile
	}
}
