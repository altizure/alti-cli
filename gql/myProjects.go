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
func MyProjects(before, after string) ([]types.Project, *types.PageInfo, error) {
	config, err := config.Load()
	if err != nil {
		return nil, nil, errors.ErrNoConfig
	}

	client := graphql.NewClient(config.Endpoint)

	// make a request
	req := graphql.NewRequest(`
		query ($before: String, $after: String) {
			my {
				projects(first: 50, before: $before, after: $after) {
					pageInfo {
						hasNextPage
						endCursor
					}
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
	req.Var("before", before)
	req.Var("after", after)

	req.Header.Set("key", config.Key)
	req.Header.Set("altitoken", config.Token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res myProjsRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, nil, err
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
	pi := types.PageInfo{
		HasNextPage:     res.My.Projects.PageInfo.HasNextPage,
		HasPreviousPage: res.My.Projects.PageInfo.HasPreviousPage,
		StartCursor:     res.My.Projects.PageInfo.StartCursor,
		EndCursor:       res.My.Projects.PageInfo.EndCursor,
	}
	return ret, &pi, nil
}

type myProjsRes struct {
	My struct {
		Projects struct {
			PageInfo struct {
				HasNextPage     bool
				HasPreviousPage bool
				StartCursor     string
				EndCursor       string
			}
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
