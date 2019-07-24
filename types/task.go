package types

import "time"

// Task represents the gql 'Task' type.
type Task struct {
	ID         string
	TaskType   string
	State      string
	StartDate  time.Time
	EndDate    time.Time
	TotalSteps int
	Step       int
	Queueing   int
}
