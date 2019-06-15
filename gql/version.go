package gql

import (
	"context"
	"time"

	"github.com/machinebox/graphql"
)

// Version gets the current version of api server.
func Version(endpoint, key string) (string, time.Duration) {
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
	start := time.Now()
	if err := client.Run(ctx, req, &res); err != nil {
		return "Offline", 0
	}
	elapsed := time.Since(start)
	return res.Versions.API, elapsed
}

type versionsRes struct {
	Versions struct {
		API string
	}
}
