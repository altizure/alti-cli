package gql

import (
	"context"
	"log"

	"github.com/jackytck/alti-cli/config"
	"github.com/machinebox/graphql"
)

// CreateProject creates a new empty project
// and returns the pid of the newly created project.
func CreateProject(name, projType, modelType, visibility string, importedModel bool) string {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($name: String!, $type: PROJECT_TYPE, $imported: Boolean, $modelType: IMPORTED_MODEL_TYPE, $visibility: PROJECT_VISIBILITY) {
			createProject(name: $name, type: $type, imported: $imported, visibility: $visibility) {
				id
			}
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// set create project variables
	req.Var("name", name)
	req.Var("type", projType)
	req.Var("imported", importedModel)
	req.Var("modelType", modelType)
	req.Var("visibility", visibility)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res createProjRes
	if err := client.Run(ctx, req, &res); err != nil {
		log.Fatal(err)
	}
	return res.CreateProject.ID
}

type createProjRes struct {
	CreateProject struct {
		ID string
	}
}
