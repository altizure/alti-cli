package cmd

import (
	"os"

	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// myprojCmd represents the myproj command
var myprojCmd = &cobra.Command{
	Use:   "myproj",
	Short: "My fist 50 projects",
	Long:  "A list of my first 50 projects.",
	Run: func(cmd *cobra.Command, args []string) {
		projs, err := gql.MyProjects()
		if err != nil {
			panic(err)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(types.ProjectHeaderString())
		for _, p := range projs {
			table.Append(p.RowString())
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(myprojCmd)
}
