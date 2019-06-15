package gql

import (
	"context"

	"github.com/machinebox/graphql"
)

// Version gets the current version of api server.
func Version(endpoint, key string) string {
	client := graphql.NewClient(endpoint + "/graphql")

	req := graphql.NewRequest(`
		{
		  versions {
		    api
		  }
		}
	`)
	req.Header.Set("key", key)

	ctx := context.Background()
	var res versionsRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "Offline"
	}
	return res.Versions.Api
}

type versionsRes struct {
	Versions struct {
		Api string
	}
}
