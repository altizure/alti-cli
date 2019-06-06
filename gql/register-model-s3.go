package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RegisterModelS3 registers a S3 model.
// And get back the registered model and the signed url to S3.
func RegisterModelS3(pid, bucket, filename string) (*types.Model, string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($id: ID!, $bucket: BucketS3Model!, $filename: String!) {
		  uploadModelS3(id: $id, bucket: $bucket, filename: $filename) {
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
	var res regModelS3Res
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, "", err
	}
	mid := res.UploadModelS3.File.ID
	url := res.UploadModelS3.URL
	if mid == "" || url == "" {
		return nil, "", errors.ErrModelReg
	}

	img := types.Model{
		ID:       res.UploadModelS3.File.ID,
		State:    res.UploadModelS3.File.State,
		Name:     res.UploadModelS3.File.Name,
		Filename: res.UploadModelS3.File.Filename,
	}

	return &img, url, nil
}

type regModelS3Res struct {
	UploadModelS3 struct {
		URL  string
		File struct {
			ID       string
			State    string
			Name     string
			Filename string
		}
	}
}
