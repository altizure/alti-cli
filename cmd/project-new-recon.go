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

// newReconCmd represents the new command
var newReconCmd = &cobra.Command{
	Use:   "recon",
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
	projNewCmd.AddCommand(newReconCmd)
	newReconCmd.Flags().StringVarP(&name, "name", "n", name, "Project name")
	newReconCmd.MarkFlagRequired("name")
	newReconCmd.Flags().StringVarP(&projType, "projectType", "p", projType, "free, pro")
	newReconCmd.Flags().StringVarP(&visibility, "visibility", "v", visibility, "public, unlisted, private")
}
