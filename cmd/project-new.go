package cmd

import (
	"fmt"
	"os"

	"github.com/jackytck/alti-cli/gql"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var name string
var projType = "free"
var visibility = "public"

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create an empty reconstruction project",
	Long:  "Create an empty reconstruction project.",
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := gql.CreateProject(name, projType, "", visibility)
		if err != nil {
			fmt.Println("Project could not be created!", err)
			return
		}

		fmt.Println("Successfully created an empty project:")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Project Type", "Visibility"})
		r := []string{pid, name, projType, visibility}
		table.Append(r)
		table.Render()
	},
}

func init() {
	projectCmd.AddCommand(newCmd)
	newCmd.Flags().StringVarP(&name, "name", "n", name, "Project name")
	newCmd.MarkFlagRequired("name")
	newCmd.Flags().StringVarP(&projType, "projectType", "p", projType, "free, pro")
	newCmd.Flags().StringVarP(&visibility, "visibility", "v", visibility, "public, unlisted, private")
}
