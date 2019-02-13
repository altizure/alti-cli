package cmd

import (
	"fmt"

	"github.com/jackytck/alti-cli/gql"
	"github.com/spf13/cobra"
)

var id string

// projRemoveCmd represents the remove command
var projRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove project by pid",
	Long:  "Remove a project by its pid.",
	Run: func(cmd *cobra.Command, args []string) {
		err := gql.RemoveProject(id)
		if err != nil {
			fmt.Println("Project could not be removed!", err)
			return
		}

		fmt.Printf("Successfully removed project with id: %q", id)
	},
}

func init() {
	projectCmd.AddCommand(projRemoveCmd)
	projRemoveCmd.Flags().StringVarP(&id, "id", "p", id, "Project id")
	projRemoveCmd.MarkFlagRequired("id")
}
