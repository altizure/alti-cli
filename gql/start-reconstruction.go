package gql

import (
	"context"
	"errors"
	"log"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// StartReconstruction starts a reconstruction by project id.
func StartReconstruction(pid string) (*types.Task, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// @TODO: start non-native task
	taskTypes, err := EnumValues("TASK_TYPE")
	if err != nil {
		return nil, err
	}
	log.Println(taskTypes)

	// make a request
	req := graphql.NewRequest(`
		mutation ($id: ID!) {
			startReconstructionWithError(id: $id) {
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
