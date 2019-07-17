package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
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
		p, err := gql.SearchProjectID(id, true)
		if err != nil {
			fmt.Println("Project could not be found! Error:", err)
			return
		}

		// check project info
		if p.NumImage > 0 || p.GigaPixel > 0 {
			var ans string
			s := ""
			if p.NumImage > 1 {
				s = "s"
			}
			fmt.Printf("Warning: Project: %q is not empty. It has %d image%s. Are you sure? (Y/N): ", p.Name, p.NumImage, s)
			fmt.Scanln(&ans)
			ans = strings.ToUpper(ans)
			if ans != "Y" && ans != service.Yes {
				fmt.Println("Cancelled.")
				return
			}
		}

		_, err = gql.RemoveProject(p.ID)
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
	projRemoveCmd.Flags().StringVarP(&id, "id", "p", id, "Project (partial) id")
	errors.Must(projRemoveCmd.MarkFlagRequired("id"))
}
