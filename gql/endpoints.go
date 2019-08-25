package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// WebEndpoint returns the current active web domain.
func WebEndpoint() string {
	config := config.Load()
	active := config.GetActive()
	ep, err := Endpoints(active.Endpoint, active.Key)
	if err != nil {
		return ""
	}
	return ep.Web
}

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
