package gql

import (
	"context"
	"net/url"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// Project return the project by the given id.
func Project(id string) (*types.Project, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		query ($id: ID!) {
			project(id: $id) {
				id
				name
				isImported
				importedState
				projectType
				numImage
				gigaPixel
				taskState
				date
				cloudPath {
					key
				}
			}
		}
	`)
	req.Var("id", id)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res projRes
	if err := client.Run(ctx, req, &res); err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, errors.ErrOffline
		default:
			return nil, err
		}
	}

	p := res.Project
	if p.ID == "" {
		return nil, errors.ErrProjNotFound
	}
	return &p, nil
}

type projRes struct {
	Project types.Project
}
