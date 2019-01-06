package cmd

import (
	"fmt"

	"github.com/jackytck/alti-cli/gql"
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
		for _, p := range projs {
			fmt.Println(p)
		}
	},
}

func init() {
	rootCmd.AddCommand(myprojCmd)
}
