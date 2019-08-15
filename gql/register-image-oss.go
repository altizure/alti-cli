package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/rand"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RefreshSTS is a HOF of GetSTS for refreshing the STS.
func RefreshSTS(pid, bucket string) func() (*types.STS, error) {
	return func() (*types.STS, error) {
		return GetSTS(pid, bucket)
	}
}

// GetSTS obtains the STS creds for this project.
func GetSTS(pid, bucket string) (*types.STS, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(`
		mutation ($pid: ID!, $bucket: BucketOSS!, $filename: String!) {
			uploadImageOSS(pid: $pid, bucket: $bucket, filename: $filename) {
				sts {
					id
					secret
					token
					bucket
					endpoint
					expire
				}
			}
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	randName, err := rand.RememberToken()
	if err != nil {
		return nil, err
	}

	// set variables
	req.Var("pid", pid)
	req.Var("bucket", bucket)
	req.Var("filename", randName)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res stsRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	id := res.UploadImageOSS.STS.ID
	if id == "" {
		return nil, errors.ErrImgReg
	}

	sts := types.STS{
		ID:       res.UploadImageOSS.STS.ID,
		Secret:   res.UploadImageOSS.STS.Secret,
		Token:    res.UploadImageOSS.STS.Token,
		Endpoint: res.UploadImageOSS.STS.Endpoint,
		Bucket:   res.UploadImageOSS.STS.Bucket,
		Expire:   res.UploadImageOSS.STS.Expire,
	}

	return &sts, nil
}

// RegisterImageOSS registers an OSS image, without getting the STS creds.
// Return the registered image.
func RegisterImageOSS(pid, bucket, filename, imageType, checksum string) (*types.Image, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($pid: ID!, $bucket: BucketOSS!, $filename: String!, $type: IMAGE_TYPE, $checksum: String) {
			uploadImageOSS(pid: $pid, bucket: $bucket, filename: $filename, type: $type, checksum: $checksum) {
				image {
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
	req.Var("type", imageType)
	req.Var("checksum", checksum)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res regImgOSSRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	iid := res.UploadImageOSS.Image.ID
	if iid == "" {
		return nil, errors.ErrImgReg
	}

	return &res.UploadImageOSS.Image, nil
}

type stsRes struct {
	UploadImageOSS struct {
		STS types.STS
	}
}

type regImgOSSRes struct {
	UploadImageOSS struct {
		Image types.Image
	}
}
