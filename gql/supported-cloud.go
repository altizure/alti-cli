package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/machinebox/graphql"
)

// SupportedCloud queries for the supported cloud of the given endpoint.
func SupportedCloud(endpoint, key string) []string {
	if endpoint == "" || key == "" {
		config := config.Load()
		active := config.GetActive()
		endpoint = active.Endpoint
		key = active.Key
	}
	client := graphql.NewClient(endpoint + "/graphql")

	req := graphql.NewRequest(`
		{
			support {
				supportedCloud
			}
		}
	`)
	req.Header.Set("key", key)

	ctx := context.Background()
	var res supCloudRes
	if err := client.Run(ctx, req, &res); err != nil {
		return []string{}
	}
	return res.Support.SupportedCloud
}

type supCloudRes struct {
	Support struct {
		SupportedCloud []string
	}
}
