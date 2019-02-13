package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/machinebox/graphql"
)

// RemoveProject removes a project by the given pid.
func RemoveProject(pid string) error {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($id: ID!) {
			removeProject(id: $id) {
				id
			}
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)
	req.Var("id", pid)

	ctx := context.Background()

	var res removeProjRes
	if err := client.Run(ctx, req, &res); err != nil {
		return err
	}
	id := res.RemoveProject.ID
	if id == "" {
		return errors.ErrProjRemove
	}
	return nil
}

type removeProjRes struct {
	RemoveProject struct {
		ID string
	}
}
