package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// StopReconstruction starts a reconstruction by project id.
func StopReconstruction(pid string) (*types.Task, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($id: ID!) {
			stopReconstruction(id: $id) {
				id
				taskType
				state
				startDate
				endDate
			}
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)
	req.Var("id", pid)

	ctx := context.Background()

	var res stopReconRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	if res.StopReconstruction.ID == "" {
		return nil, errors.ErrTaskStop
	}
	return &res.StopReconstruction, nil
}

type stopReconRes struct {
	StopReconstruction types.Task
}
