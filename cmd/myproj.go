package cmd

import (
	"fmt"
	"os"

	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/types"
	"github.com/spf13/cobra"
)

// myprojCmd represents the myproj command
var myprojCmd = &cobra.Command{
	Use:   "myproj",
	Short: "My fist 50 projects",
	Long:  "A list of my first 50 projects.",
	Run: func(cmd *cobra.Command, args []string) {
		projs, page, err := gql.MyProjects("", "")
		if err != nil {
			panic(err)
		}
		table := types.ProjectsToTable(projs, os.Stdout)
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(myprojCmd)
}
