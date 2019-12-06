package cmd

import (
	"log"
	"os"

	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// taskTypeCmd represents the taskType command
var taskTypeCmd = &cobra.Command{
	Use:   "task-type",
	Short: "List all available task types",
	Long:  `'alti-cli list task-type' to list all available task types used in the 'project start' command.`,
	Run: func(cmd *cobra.Command, args []string) {
		// pre-checks
		if err := service.Check(
			nil,
			service.CheckAPIServer(),
		); err != nil {
			log.Println(err)
			return
		}

		tts, err := gql.EnumValues("TASK_TYPE")
		if err != nil {
			panic(err)
		}
		var rows [][]string
		for _, t := range tts {
			rows = append(rows, []string{t})
		}

		// render
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Task Type"})
		table.AppendBulk(rows)
		table.Render()
	},
}

func init() {
	listCmd.AddCommand(taskTypeCmd)
}
