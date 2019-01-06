package types

import "fmt"

// Project represents the gql project type.
type Project struct {
	ID          string
	Name        string
	IsImported  bool
	ProjectType string
	NumImage    int
	GigaPixel   float64
	TaskState   string
}

func (p Project) String() string {
	pat := "ID: %s\tName: %s\tIsImported: %v\tProjectType: %s\tNumImage: %d\tGigaPixel: %.2f\tTaskState: %s"
	return fmt.Sprintf(pat, p.ID, p.Name, p.IsImported, p.ProjectType, p.NumImage, p.GigaPixel, p.TaskState)
}
