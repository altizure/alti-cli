package gql

import (
	"context"
	"net/url"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
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
		switch err.(type) {
		case *url.Error:
			return nil, nil, 0, errors.ErrOffline
		default:
			return nil, nil, 0, err
		}
	}

	var ret []types.Project
	for _, e := range res.My.AllProjects.Edges {
		ret = append(ret, e.Node)
	}
	pi := res.My.AllProjects.PageInfo
	return ret, &pi, res.My.AllProjects.TotalCount, nil
}

type myProjsRes struct {
	My struct {
		AllProjects struct {
			TotalCount int
			PageInfo   types.PageInfo
			Edges      []struct {
				Node types.Project
			}
		}
	}
}
