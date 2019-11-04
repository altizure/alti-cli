package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// ReportProject reports a project with error description.
func ReportProject(pid, desc string) error {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(`
		mutation ($pid: ID!, $desc: String!) {
			reportProject(id: $pid, description: $desc) {
				id
			}
		}
	`)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	req.Var("pid", pid)
	req.Var("desc", desc)

	ctx := context.Background()
	var res reportProjRes
	if err := client.Run(ctx, req, &res); err != nil {
		return err
	}
	if res.ReportProject.ID == "" {
		return errors.ErrReportProj
	}
	return nil
}

type reportProjRes struct {
	ReportProject types.Project
}
