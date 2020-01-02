package gql

import (
	"context"
	"net/url"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// AllProjectImages queries all of the project images by cursor.
func AllProjectImages(pid string, first, last int, before, after string) ([]types.ProjectImage, *types.PageInfo, int, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		query ($id: ID!, $first: Int, $last: Int, $before: String, $after: String) {
			project(id: $id) {
				allImages(first: $first, last: $last, before: $before, after: $after) {
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
							state
							grounded
							url
						}
					}
				}
			}
		}
	`)
	req.Var("id", pid)
	if first > 0 {
		req.Var("first", first)
	}
	if last > 0 {
		req.Var("last", last)
	}
	req.Var("before", before)
	req.Var("after", after)

	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res allImgsRes
	if err := client.Run(ctx, req, &res); err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, nil, 0, errors.ErrOffline
		default:
			return nil, nil, 0, err
		}
	}

	var ret []types.ProjectImage
	for _, e := range res.Project.AllImages.Edges {
		ret = append(ret, e.Node)
	}
	pi := res.Project.AllImages.PageInfo
	return ret, &pi, res.Project.AllImages.TotalCount, nil
}

type allImgsRes struct {
	Project struct {
		AllImages struct {
			TotalCount int
			PageInfo   types.PageInfo
			Edges      []struct {
				Node types.ProjectImage
			}
		}
	}
}
