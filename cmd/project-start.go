package cmd

import (
	"fmt"
	"os"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// startReconCmd represents the start reconstruction command
var startReconCmd = &cobra.Command{
	Use:   "start",
	Short: "Start reconstruction.",
	Long:  "Start a native reconstruction of a project.",
	Run: func(cmd *cobra.Command, args []string) {
		p, err := gql.SearchProjectID(id, true)
		if err != nil {
			fmt.Println("Project could not be found! Error:", err)
			return
		}

		t, err := gql.StartReconstruction(p.ID)
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

		fmt.Printf("Successfully started a %q task with state: %q\n", t.TaskType, t.State)
	},
}

func init() {
	projectCmd.AddCommand(startReconCmd)
	startReconCmd.Flags().StringVarP(&id, "id", "p", id, "Project (partial) id")
	errors.Must(startReconCmd.MarkFlagRequired("id"))
}
