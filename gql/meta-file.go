package gql

import (
	"context"
	"net/url"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// ProjectMetaFile return the info of a project meta file.
func ProjectMetaFile(pid, iid string) (*types.MetaFile, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		query ($pid: ID!, $iid: ID!) {
			project(id: $pid) {
				metaFile(id: $iid) {
					id
					state
					name
					filename
					filesize
					date
					checksum
					error
				}
			}
		}
	`)
	req.Var("pid", pid)
	req.Var("iid", iid)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res projMetaRes
	if err := client.Run(ctx, req, &res); err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, errors.ErrOffline
		default:
			return nil, err
		}
	}

	m := res.Project.MetaFile
	if m.State == "" {
		return nil, errors.ErrMetaNotFound
	}
	return &m, nil
}

type projMetaRes struct {
	Project struct {
		MetaFile types.MetaFile
	}
}
