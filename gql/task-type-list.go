package gql

import (
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/text"
)

// QueryTaskType infers the exact task type name from query string taskType.
func QueryTaskType(taskType string) (string, []string, error) {
	list, err := EnumValues("TASK_TYPE")
	if err != nil {
		return "", list, err
	}
	ret := text.BestMatch(list, taskType, "")
	if ret == "" {
		return ret, list, errors.ErrTaskTypeInvalid
	}
	return ret, list, nil
}
