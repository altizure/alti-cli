package cmd

import (
	"fmt"
	"os"
	"strings"

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
		// check project info
		p, err := gql.Project(id)
		if err != nil {
			fmt.Println("Project could not be removed! Error:", err)
			return
		}
		if p.NumImage > 0 || p.GigaPixel > 0 {
			var ans string
			fmt.Printf("Warning: Project: %q is not empty. Are you sure? (Y/N): ", p.Name)
			fmt.Scanln(&ans)
			ans = strings.ToUpper(ans)
			if ans != "Y" && ans != "YES" {
				fmt.Println("Cancelled.")
				return
			}
		}

		_, err = gql.RemoveProject(id)
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
