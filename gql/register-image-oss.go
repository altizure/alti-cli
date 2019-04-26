package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RegisterImageOSS registers an OSS image, without getting the STS creds.
// Return the registerd image.
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

	img := types.Image{
		ID:       res.UploadImageOSS.Image.ID,
		State:    res.UploadImageOSS.Image.State,
		Name:     res.UploadImageOSS.Image.Name,
		Filename: res.UploadImageOSS.Image.Filename,
	}

	return &img, nil
}

type regImgOSSRes struct {
	UploadImageOSS struct {
		Image struct {
			ID       string
			State    string
			Name     string
			Filename string
		}
	}
}
