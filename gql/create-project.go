package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/machinebox/graphql"
)

// CreateProject creates a new empty project
// and returns the pid of the newly created project.
func CreateProject(name, projType, modelType, visibility string) (string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($name: String!, $type: PROJECT_TYPE, $imported: Boolean, $modelType: IMPORTED_MODEL_TYPE, $visibility: PROJECT_VISIBILITY) {
			createProject(name: $name, type: $type, imported: $imported, modelType: $modelType, visibility: $visibility) {
				id
			}
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// set create project variables
	req.Var("name", name)
	req.Var("type", projType)
	if modelType != "" {
		req.Var("modelType", modelType)
		req.Var("imported", true)
	}
	req.Var("visibility", visibility)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res createProjRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	pid := res.CreateProject.ID
	if pid == "" {
		return "", errors.ErrProjCreate
	}
	return pid, nil
}

type createProjRes struct {
	CreateProject struct {
		ID string
	}
}
