package gql

import (
	"fmt"

	"github.com/jackytck/alti-cli/config"
	"github.com/machinebox/graphql"
)

// ActiveClient constructs the gql client for the currently active profile.
// Return the gql client, endpint, key and token.
func ActiveClient(room string) (*graphql.Client, string, string, string) {
	if room == "" {
		room = "graphql"
	}

	config := config.Load()
	active := config.GetActive()
	endpoint := active.Endpoint
	key := active.Key
	token := active.Token

	url := fmt.Sprintf("%s/%s", endpoint, room)
	client := graphql.NewClient(url)

	return client, endpoint, key, token
}
