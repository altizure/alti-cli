package cmd

import (
	"fmt"
	"os"

	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/types"
	"github.com/spf13/cobra"
)

var search string
var pageCount = 12

// myprojCmd represents the myproj command
var myprojCmd = &cobra.Command{
	Use:   "myproj",
	Short: "My fist 50 projects",
	Long:  "A list of my first 50 projects.",
	Run: func(cmd *cobra.Command, args []string) {
		projs, page, total, err := gql.MyProjects(pageCount, 0, "", "", search)
		if err != nil {
			panic(err)
		}
		table := types.ProjectsToTable(projs, os.Stdout)
		table.Render()
		fmt.Printf("Total: %d\n", total)
		if page.HasNextPage {
			fmt.Println("More projects by: 'alti-cli myproj more'")
		}
	},
}

func init() {
	rootCmd.AddCommand(myprojCmd)
	myprojCmd.Flags().IntVarP(&pageCount, "count", "c", pageCount, "number of projects to fetch")
	myprojCmd.Flags().StringVarP(&search, "search", "q", search, "display name to search")
}
