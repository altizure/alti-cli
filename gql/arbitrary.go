package gql

import (
	"context"
	"encoding/json"

	"github.com/jackytck/alti-cli/config"
	"github.com/machinebox/graphql"
)

// Arbitrary makes arbitrary query or mutation.
func Arbitrary(query string) (string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(query)

	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	ctx := context.Background()
	var res json.RawMessage
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	return PrettyPrint(res)
}
