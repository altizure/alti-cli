package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/machinebox/graphql"
)

// HasMetaFile asks if the project has the given meta file by hash.
func HasMetaFile(pid, checksum string) (bool, error) {
	if pid == "" || checksum == "" {
		return false, nil
	}

	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		query ($pid: ID!, $checksum: String!) {
			project(id: $pid) {
				hasMetaFile(checksum: $checksum) {
					id
					state
				}
			}
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)
	req.Var("pid", pid)
	req.Var("checksum", checksum)

	ctx := context.Background()

	var res hasMetaRes
	if err := client.Run(ctx, req, &res); err != nil {
		return false, err
	}
	id := res.Project.HasMetaFile.ID
	if id == "" {
		return false, nil
	}
	return true, nil
}

type hasMetaRes struct {
	Project struct {
		HasMetaFile struct {
			ID    string
			State string
		}
	}
}
