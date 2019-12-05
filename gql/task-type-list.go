package gql

import (
	"context"
	"sort"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/text"
	"github.com/machinebox/graphql"
)

// QueryTaskType infers the exact task type name from query string taskType.
func QueryTaskType(taskType string) (string, []string, error) {
	list, err := TaskTypeList()
	if err != nil {
		return "", list, err
	}
	ret := text.BestMatch(list, taskType, "")
	if ret == "" {
		return ret, list, errors.ErrTaskTypeInvalid
	}
	return ret, list, nil
}

// TaskTypeList returns a list of available task types supported by the api server.
func TaskTypeList() ([]string, error) {
	var ret []string

	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(`
		query ($type: String!) {
			__type(name: $type) {
				enumValues {
					name
				}
			}
		}
	`)

	req.Header.Set("key", active.Key)
	req.Var("type", "TASK_TYPE")

	ctx := context.Background()
	var res taskTypeRes
	if err := client.Run(ctx, req, &res); err != nil {
		return ret, err
	}

	for _, c := range res.Type.EnumValues {
		ret = append(ret, c.Name)
	}
	sort.Strings(ret)

	return ret, nil
}

type taskTypeRes struct {
	Type enumTTType `json:"__type"`
}

type enumTTType struct {
	EnumValues []struct {
		Name string
	}
}
