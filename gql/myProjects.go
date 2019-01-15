package gql

import (
	"context"
	"time"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// MyProjects queries simple info of my first 50 projects.
func MyProjects(first, last int, before, after, search string) ([]types.Project, *types.PageInfo, int, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		query ($first: Int, $last: Int, $before: String, $after: String, $search: String) {
			my {
				allProjects(first: $first, last: $last, before: $before, after: $after, search: $search) {
					totalCount
					pageInfo {
						hasPreviousPage
						hasNextPage
						startCursor
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
	if first > 0 {
		req.Var("first", first)
	}
	if last > 0 {
		req.Var("last", last)
	}
	req.Var("before", before)
	req.Var("after", after)
	req.Var("search", search)

	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res myProjsRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, nil, 0, err
	}

	var ret []types.Project
	for _, e := range res.My.AllProjects.Edges {
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
		HasNextPage:     res.My.AllProjects.PageInfo.HasNextPage,
		HasPreviousPage: res.My.AllProjects.PageInfo.HasPreviousPage,
		StartCursor:     res.My.AllProjects.PageInfo.StartCursor,
		EndCursor:       res.My.AllProjects.PageInfo.EndCursor,
	}
	return ret, &pi, res.My.AllProjects.TotalCount, nil
}

type myProjsRes struct {
	My struct {
		AllProjects struct {
			TotalCount int
			PageInfo   struct {
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
