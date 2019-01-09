package cmd

import (
	"os"
	"time"

	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/types"
	"github.com/spf13/cobra"
)

// moreCmd represents the more command
var moreCmd = &cobra.Command{
	Use:   "more",
	Short: "List of all of my projects paginated.",
	Long:  "Show all of my projects in pages. @TODO: Go to the next page by pressing n or Space or Enter. Previous page by p or Shift+Space or Shift+Enter.",
	Run: func(cmd *cobra.Command, args []string) {
		after := ""
		for {
			projs, page, err := gql.MyProjects("", after)
			if err != nil {
				panic(err)
			}
			table := types.ProjectsToTable(projs, os.Stdout)
			table.Render()
			if !page.HasNextPage {
				break
			}
			after = page.EndCursor
			time.Sleep(time.Second * 3)
		}
	},
}

func init() {
	myprojCmd.AddCommand(moreCmd)
}
