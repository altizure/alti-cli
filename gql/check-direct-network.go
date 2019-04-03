package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/machinebox/graphql"
)

// CheckDirectNetwork tests if the api server could reach this client.
func CheckDirectNetwork(url string) bool {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(`
		query ($url: String!) {
			support {
				networkTest(url: $url)
			}
		}
	`)
	req.Var("url", url)

	req.Header.Set("key", active.Key)

	ctx := context.Background()
	var res networkTestRes
	if err := client.Run(ctx, req, &res); err != nil {
		return false
	}
	success := res.Support.NetworkTest
	return success == "Success"
}

type networkTestRes struct {
	Support struct {
		NetworkTest string
	}
}
