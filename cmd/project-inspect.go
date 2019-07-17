package cmd

import (
	"fmt"
	"os"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/types"
	"github.com/spf13/cobra"
)

// inspectCmd represents the inspect command
var projInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect any individual project by partial id.",
	Long:  "Inspect any individual project details given by its id. ID could be a partial id.",
	Run: func(cmd *cobra.Command, args []string) {
		p, err := gql.SearchProjectID(id, false)
		if err != nil {
			fmt.Println("Project could not be found! Error:", err)
			return
		}
		table := types.ProjectsToTable([]types.Project{*p}, os.Stdout)
		table.Render()
	},
}

func init() {
	projectCmd.AddCommand(projInspectCmd)
	projInspectCmd.Flags().StringVarP(&id, "id", "p", id, "(Partial) Project id")
	errors.Must(projInspectCmd.MarkFlagRequired("id"))
}
