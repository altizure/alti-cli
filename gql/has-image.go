package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/machinebox/graphql"
)

// HasImage asks if the project has the given image by hash.
func HasImage(pid, checksum string) (bool, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation hasImage($pid: ID!, $checksum: String!) {
		  hasImage(pid: $pid, checksum: $checksum) {
		    id
		    state
		  }
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)
	req.Var("id", pid)
	req.Var("checksum", checksum)

	ctx := context.Background()

	var res hasImageRes
	if err := client.Run(ctx, req, &res); err != nil {
		return false, err
	}
	id := res.HasImage.ID
	if id == "" {
		return false, nil
	}
	return true, nil
}

type hasImageRes struct {
	HasImage struct {
		ID    string
		State string
	}
}
