package gql

import (
	"context"
	"time"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// RemoveProject removes a project by the given pid.
func RemoveProject(pid string) (*types.Project, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation ($id: ID!) {
			removeProject(id: $id) {
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
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)
	req.Var("id", pid)

	ctx := context.Background()

	var res removeProjRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	id := res.RemoveProject.ID
	if id == "" {
		return nil, errors.ErrProjRemove
	}
	return res.toProject(), nil
}

type removeProjRes struct {
	RemoveProject struct {
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

func (rp *removeProjRes) toProject() *types.Project {
	r := rp.RemoveProject
	return &types.Project{
		ID:          r.ID,
		Name:        r.Name,
		IsImported:  r.IsImported,
		ProjectType: r.ProjectType,
		NumImage:    r.NumImage,
		GigaPixel:   r.GigaPixel,
		TaskState:   r.TaskState,
		Date:        r.Date,
	}
}
