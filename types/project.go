package types

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jackytck/alti-cli/config"
	"github.com/olekukonko/tablewriter"
)

// Project represents the gql project type.
type Project struct {
	ID            string
	Name          string
	IsImported    bool
	ImportedState string
	ProjectType   string
	NumImage      int
	GigaPixel     float64
	TaskState     string
	Date          time.Time
	CloudPath     []CloudPath
	Downloads     DownloadsConnection
}

func (p Project) String() string {
	pat := "ID: %s\tName: %s\tIsImported: %v\tProjectType: %s\tNumImage: %d\tGigaPixel: %.2f\tTaskState: %s\tCloud: %v"
	return fmt.Sprintf(pat, p.ID, p.Name, p.IsImported, p.ProjectType, p.NumImage, p.GigaPixel, p.TaskState, strings.Join(p.Cloud(), ", "))
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
		"Cloud",
		"Date",
		"Model Link",
	}
}

// Cloud returns the cloud keys of the project.
func (p Project) Cloud() []string {
	var ret []string
	for _, c := range p.CloudPath {
		ret = append(ret, c.Key)
	}
	return ret
}

// RowString gives a row of string for the table output.
func (p Project) RowString(endPoint string) []string {
	base := ""
	if strings.Contains(endPoint, "altizure") {
		base = strings.Replace(endPoint, "api", "www", 1)
	} else if strings.Contains(endPoint, "8082") {
		base = strings.Replace(endPoint, "8082", "8091", 1)
	}
	return []string{
		p.ID,
		p.Name,
		fmt.Sprintf("%v", p.IsImported),
		p.ProjectType,
		fmt.Sprintf("%d", p.NumImage),
		fmt.Sprintf("%.2f", p.GigaPixel),
		p.TaskState,
		strings.Join(p.Cloud(), ", "),
		p.Date.Format("2006-01-02 15:04:05"),
		fmt.Sprintf("%s/project-model?pid=%v", base, p.ID),
	}
}

// ProjectsToTable transforms slice of projects into a table.
func ProjectsToTable(ps []Project, w io.Writer) *tablewriter.Table {
	config := config.Load()
	active := config.GetActive()

	table := tablewriter.NewWriter(w)
	table.SetHeader(ProjectHeaderString())
	for _, p := range ps {
		table.Append(p.RowString(active.Endpoint))
	}
	return table
}
