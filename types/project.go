package types

import (
	"fmt"
	"time"
)

// Project represents the gql project type.
type Project struct {
	ID          string
	Name        string
	IsImported  bool
	ProjectType string
	NumImage    int
	GigaPixel   float64
	TaskState   string
	Date        time.Time
}

func (p Project) String() string {
	pat := "ID: %s\tName: %s\tIsImported: %v\tProjectType: %s\tNumImage: %d\tGigaPixel: %.2f\tTaskState: %s"
	return fmt.Sprintf(pat, p.ID, p.Name, p.IsImported, p.ProjectType, p.NumImage, p.GigaPixel, p.TaskState)
}

// ProjectHeaderString gives a row of string for the table header.
func ProjectHeaderString() []string {
	return []string{
		"ID",
		"Name",
		"Is Imported",
		"Project Type",
		"Num Image",
		"Giga-Pixel",
		"Task State",
		"Date",
	}
}

// RowString gives a row of string for the table output.
func (p Project) RowString() []string {
	return []string{
		p.ID,
		p.Name,
		fmt.Sprintf("%v", p.IsImported),
		p.ProjectType,
		fmt.Sprintf("%d", p.NumImage),
		fmt.Sprintf("%.2f", p.GigaPixel),
		p.TaskState,
		p.Date.Format("2006-01-02 15:04:05"),
	}
}
