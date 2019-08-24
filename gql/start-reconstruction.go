package gql

import (
	"context"
	"errors"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// StartReconstruction starts a reconstruction by project id.
func StartReconstruction(pid, taskType string) (*types.Task, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($id: ID!, $taskType: TASK_TYPE) {
			startReconstructionWithError(id: $id, options: {taskType: $taskType}) {
				error {
					code
					message
				}
				task {
					id
					taskType
					state
					startDate
					queueing
				}
			}
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)
	req.Var("id", pid)
	req.Var("taskType", taskType)

	ctx := context.Background()

	var res startReconRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	err2 := res.StartReconstructionWithError.Error
	if err2.Message != "" {
		return nil, errors.New(err2.Message)
	}
	return &res.StartReconstructionWithError.Task, nil
}

type startReconRes struct {
	StartReconstructionWithError struct {
		Error types.Error
		Task  types.Task
	}
}
