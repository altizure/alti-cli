package gql

import (
	"context"
	"net/url"
	"time"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// Project return the project by the given id.
func Project(id string) (*types.Project, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		query ($id: ID!) {
			project(id: $id) {
				id
				name
				isImported
				projectType
				numImage
				gigaPixel
				taskState
				date
			}
		}
	`)
	req.Var("id", id)
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res projRes
	if err := client.Run(ctx, req, &res); err != nil {
		switch err.(type) {
		case *url.Error:
			return nil, errors.ErrOffline
		default:
			return nil, err
		}
	}

	r := res.Project
	if r.ID == "" {
		return nil, errors.ErrProjNotFound
	}

	p := types.Project{
		ID:          r.ID,
		Name:        r.Name,
		IsImported:  r.IsImported,
		ProjectType: r.ProjectType,
		NumImage:    r.NumImage,
		GigaPixel:   r.GigaPixel,
		TaskState:   r.TaskState,
		Date:        r.Date,
	}

	return &p, nil
}

type projRes struct {
	Project struct {
		ID          string
		Name        string
		IsImported  bool
		ProjectType string
		NumImage    int
		GigaPixel   float64
		TaskState   string
		Date        time.Time
	}
}
