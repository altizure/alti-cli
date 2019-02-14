package cmd

import (
	"fmt"
	"os"

	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/types"
	"github.com/spf13/cobra"
)

var id string

// projRemoveCmd represents the remove command
var projRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove project by pid",
	Long:  "Remove a project by its pid.",
	Run: func(cmd *cobra.Command, args []string) {
		p, err := gql.RemoveProject(id)
		if err != nil {
			fmt.Println("Project could not be removed! Error:", err)
			return
		}

		fmt.Println("Successfully removed project:")
		table := types.ProjectsToTable([]types.Project{*p}, os.Stdout)
		table.Render()
	},
}

func init() {
	projectCmd.AddCommand(projRemoveCmd)
	projRemoveCmd.Flags().StringVarP(&id, "id", "p", id, "Project id")
	projRemoveCmd.MarkFlagRequired("id")
}
