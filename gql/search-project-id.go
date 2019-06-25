package gql

import (
	"context"
	"net/url"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// SearchProjectID returns the latest project id by the given partial id.
func SearchProjectID(id string, myProj bool) (*types.Project, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		query ($id: String, $limit:Int, $myProj: Boolean) {
		  search {
		    projectID(id: $id, limit: $limit, myProject: $myProj) {
		      id
					name
					isImported
					projectType
					numImage
					gigaPixel
					taskState
					date
		    }
		  }
		}
	`)
	req.Var("id", id)
	req.Var("limit", 1)
	req.Var("myProj", myProj)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// define a Context for the request
	ctx := context.Background()

	var res searchProjRes
	if err := client.Run(ctx, req, &res); err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, errors.ErrOffline
		default:
			return nil, err
		}
	}

	ps := res.Search.ProjectID
	if len(ps) == 0 {
		return nil, errors.ErrProjNotFound
	}

	return &ps[0], nil
}

type searchProjRes struct {
	Search struct {
		ProjectID []types.Project
	}
}
