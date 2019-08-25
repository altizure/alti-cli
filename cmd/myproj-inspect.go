package cmd

import (
	"fmt"
	"os"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/types"
	"github.com/spf13/cobra"
)

// myprojInspectCmd represents the myprojInspect command
var myprojInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect my individual project by partial id.",
	Long:  "Inspect my individual project details given by its id. ID could be a partial id.",
	Run: func(cmd *cobra.Command, args []string) {
		p, err := gql.SearchProjectID(id, true)
		if err != nil {
			fmt.Println("Project could not be found! Error:", err)
			return
		}
		table := types.ProjectsToTable([]types.Project{*p}, gql.WebEndpoint(), os.Stdout)
		table.Render()
	},
}

func init() {
	myprojCmd.AddCommand(myprojInspectCmd)
	myprojInspectCmd.Flags().StringVarP(&id, "id", "p", id, "(Partial) Project id")
	errors.Must(myprojInspectCmd.MarkFlagRequired("id"))
}
