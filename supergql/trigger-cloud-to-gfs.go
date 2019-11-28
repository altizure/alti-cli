package supergql

import (
	"context"
	"errors"
	"strings"
)

// TriggerCloudToGFS triggers an op to download images from cloud to gfs.
func TriggerCloudToGFS(pid string) (string, error) {
	gql := `
		query ($pid: ID!) {
			triggerCloudToGFS(id: $pid) {
				error {
					message
				}
				state
			}
		}
	`
	req, client := SuperRequest(gql)
	req.Var("pid", pid)

	ctx := context.Background()

	var res tcgRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	if res.TriggerCloudToGFS.Error.Message != "" {
		return "", errors.New(strings.ToLower(res.TriggerCloudToGFS.Error.Message))
	}

	return res.TriggerCloudToGFS.State, nil
}

type tcgRes struct {
	TriggerCloudToGFS struct {
		Error struct {
			Message string
		}
		State string
	}
}
