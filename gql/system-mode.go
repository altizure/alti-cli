package gql

import (
	"context"

	"github.com/machinebox/graphql"
)

// CheckSystemMode checks if the api server is in Normal, ReadOnly or Offline mode.
func CheckSystemMode(endpoint, key string) string {
	client := graphql.NewClient(endpoint + "/graphql")

	req := graphql.NewRequest(`
		{
		  support {
		    systemMode
		  }
		}
	`)
	req.Header.Set("key", key)

	ctx := context.Background()
	var res systemModeRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "Offline"
	}
	mode := res.Support.SystemMode
	if mode == "" {
		mode = "Forbidden"
	}
	return mode
}

// ActiveSystemMode checks the system mode of currently active profile.
func ActiveSystemMode() string {
	_, endpoint, key, _ := ActiveClient("")
	return CheckSystemMode(endpoint, key)
}

type systemModeRes struct {
	Support struct {
		SystemMode string
	}
}
