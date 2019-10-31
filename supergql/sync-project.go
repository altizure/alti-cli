package supergql

import (
	"context"
	"errors"
	"strings"
)

// SyncProject sync a project and return its progress state.
func SyncProject(pid string) (string, error) {
	gql := `
		mutation ($pid: ID!) {
			triggerCloudSync(id: $pid) {
				error {
					message
				}
				progress {
					state
				}
			}
		}
	`
	req, client := SuperRequest(gql)
	req.Var("pid", pid)

	ctx := context.Background()

	var res syncRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	if res.TriggerCloudSync.Error.Message != "" {
		return "", errors.New(strings.ToLower(res.TriggerCloudSync.Error.Message))
	}

	return res.TriggerCloudSync.Progress.State, nil
}

type syncRes struct {
	TriggerCloudSync struct {
		Error struct {
			Message string
		}
		Progress struct {
			State string
		}
	}
}
