package gql

import (
	"context"

	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// Endpoints gets the endpoints of altizure servers.
func Endpoints(endpoint, key string) (*types.Endpoints, error) {
	client := graphql.NewClient(endpoint + "/graphql")

	req := graphql.NewRequest(`
		{
			support {
				endpoints {
					api
					web
				}
			}
		}
	`)
	req.Header.Set("key", key)

	ctx := context.Background()
	var res endpointsRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res.Support.Endpoints, nil
}

type endpointsRes struct {
	Support struct {
		Endpoints types.Endpoints
	}
}
