package gql

import (
	"context"
	"net/url"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// ProjectImage return the info of a project image.
func ProjectImage(pid, iid string) (*types.Image, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		query ($pid: ID!, $iid: ID!) {
			project(id: $pid) {
				image(id: $iid) {
					id
					state
					name
					filename
					error
				}
			}
		}
	`)
	req.Var("pid", pid)
	req.Var("iid", iid)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res projImgRes
	if err := client.Run(ctx, req, &res); err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, errors.ErrOffline
		default:
			return nil, err
		}
	}

	i := res.Project.Image
	if i.State == "" {
		return nil, errors.ErrImgNotFound
	}
	return &i, nil
}

type projImgRes struct {
	Project struct {
		Image types.Image
	}
}
