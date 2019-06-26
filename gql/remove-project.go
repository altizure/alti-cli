package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RemoveProject removes a project by the given pid.
func RemoveProject(pid string) (*types.Project, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($id: ID!) {
			removeProject(id: $id) {
				id
				name
				isImported
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
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)
	req.Var("id", pid)

	ctx := context.Background()

	var res removeProjRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	id := res.RemoveProject.ID
	if id == "" {
		return nil, errors.ErrProjRemove
	}
	return &res.RemoveProject, nil
}

type removeProjRes struct {
	RemoveProject types.Project
}
