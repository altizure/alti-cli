package gql

import (
	"context"
	"time"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// MyProjects queries simple info of my first 50 projects.
func MyProjects() ([]types.Project, error) {
	config, err := config.Load()
	if err != nil {
		return nil, errors.ErrNoConfig
	}

	client := graphql.NewClient(config.Endpoint)

	// make a request
	req := graphql.NewRequest(`
		query {
			my {
				projects(first: 50) {
					edges {
						node {
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
			}
		}
	`)
	req.Header.Set("key", config.Key)
	req.Header.Set("altitoken", config.Token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res myProjsRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}

	var ret []types.Project
	for _, e := range res.My.Projects.Edges {
		n := e.Node
		p := types.Project{
			ID:          n.ID,
			Name:        n.Name,
			IsImported:  n.IsImported,
			ProjectType: n.ProjectType,
			NumImage:    n.NumImage,
			GigaPixel:   n.GigaPixel,
			TaskState:   n.TaskState,
			Date:        n.Date,
		}
		ret = append(ret, p)
	}
	return ret, nil
}

type myProjsRes struct {
	My struct {
		Projects struct {
			Edges []struct {
				Node struct {
					ID          string
					Name        string
					IsImported  bool
					ProjectType string
					NumImage    int
					GigaPixel   float64
					TaskState   string
					Date        time.Time
				}
			}
		}
	}
}
