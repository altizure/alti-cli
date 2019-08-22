package cmd

import (
	"fmt"
	"os"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// stopReconCmd represents the stop reconstruction command
var stopReconCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop reconstruction.",
	Long:  "Stop the most current task (usually a native reconstruction) of a project.",
	Run: func(cmd *cobra.Command, args []string) {
		p, err := gql.SearchProjectID(id, true)
		if err != nil {
			fmt.Println("Project could not be found! Error:", err)
			return
		}

		t, err := gql.StopReconstruction(p.ID)
		if err != nil {
			fmt.Printf("Error: %q\n", err.Error())
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Task Type", "State", "Start Date", "Queueing"})
		d := t.StartDate.Format("2006-01-02 15:04:05")
		r := []string{t.ID, t.TaskType, t.State, d, string(t.Queueing)}
		table.Append(r)
		table.Render()

		fmt.Printf("Successfully stopped a %q task with state: %q\n", t.TaskType, t.State)
	},
}

func init() {
	projectCmd.AddCommand(stopReconCmd)
	stopReconCmd.Flags().StringVarP(&id, "id", "p", id, "Project (partial) id")
	errors.Must(stopReconCmd.MarkFlagRequired("id"))
}
