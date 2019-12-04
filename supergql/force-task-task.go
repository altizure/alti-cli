package supergql

import (
	"context"
	"errors"
	"strings"
)

// ForceStartTask forces start any task.
func ForceStartTask(pid, taskType string) (string, string, error) {
	gql := `
		mutation ($pid: ID!) {
			forceStartTask(id: $pid) {
				error {
					message
				}
				task {
					id
					state
				}
			}
		}
	`
	if taskType != "" {
		gql = `
			mutation ($pid: ID!, $taskType: TASK_TYPE) {
				forceStartTask(id: $pid, options: {taskType: $taskType}) {
					error {
						message
					}
					task {
						id
						state
					}
				}
			}
		`
	}

	req, client := SuperRequest(gql)
	req.Var("pid", pid)
	if taskType != "" {
		req.Var("taskType", taskType)
	}
	ctx := context.Background()

	var res forceStartRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", "", err
	}

	e := res.ForceStartTask.Error
	t := res.ForceStartTask.Task
	if e.Message != "" {
		return "", "", errors.New(strings.ToLower(e.Message))
	}

	return t.ID, t.State, nil
}

type forceStartRes struct {
	ForceStartTask struct {
		Error struct {
			Message string
		}
		Task struct {
			ID    string
			State string
		}
	}
}
