package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RegisterImageMinio registers a minio image.
// And get back the registered image and the signed url to minio.
func RegisterImageMinio(pid, bucket, filename, imageType, checksum string) (*types.Image, string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($pid: ID!, $bucket: BucketMinio!, $filename: String!, $type: IMAGE_TYPE, $checksum: String) {
			uploadImageMinio(pid: $pid, bucket: $bucket, filename: $filename, type: $type, checksum: $checksum) {
				url
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
	var res regImgMinioRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, "", err
	}
	iid := res.UploadImageMinio.Image.ID
	url := res.UploadImageMinio.URL
	if iid == "" || url == "" {
		return nil, "", errors.ErrImgReg
	}

	return &res.UploadImageMinio.Image, url, nil
}

type regImgMinioRes struct {
	UploadImageMinio struct {
		URL   string
		Image types.Image
	}
}
