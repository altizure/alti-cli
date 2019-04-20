package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RegisterImageS3 registers a S3 image.
// And get back the registered image and the signed url to S3.
func RegisterImageS3(pid, bucket, filename, imageType, checksum string) (*types.Image, string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($pid: ID!, $bucket: BucketS3!, $filename: String!, $type: IMAGE_TYPE, $checksum: String) {
			uploadImageS3(pid: $pid, bucket: $bucket, filename: $filename, type: $type, checksum: $checksum) {
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
	var res regImgS3Res
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, "", err
	}
	iid := res.UploadImageS3.Image.ID
	url := res.UploadImageS3.URL
	if iid == "" || url == "" {
		return nil, "", errors.ErrImgReg
	}

	img := types.Image{
		ID:       res.UploadImageS3.Image.ID,
		State:    res.UploadImageS3.Image.State,
		Name:     res.UploadImageS3.Image.Name,
		Filename: res.UploadImageS3.Image.Filename,
	}

	return &img, url, nil
}

type regImgS3Res struct {
	UploadImageS3 struct {
		URL   string
		Image struct {
			ID       string
			State    string
			Name     string
			Filename string
		}
	}
}
