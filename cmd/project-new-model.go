package cmd

import (
	"fmt"
	"os"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var modelType = "CAD"

// newModelCmd represents the new command
var newModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Create an empty model project",
	Long:  "Create an empty model project.",
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := gql.CreateProject(name, projType, modelType, visibility)
		if err != nil {
			fmt.Println("Project could not be created!", err)
			return
		}

		fmt.Println("Successfully created an empty imported project:")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Project Type", "ModelType", "Visibility"})
		r := []string{pid, name, projType, modelType, visibility}
		table.Append(r)
		table.Render()

		newPID = pid
	},
}

func init() {
	projNewCmd.AddCommand(newModelCmd)
	newModelCmd.Flags().StringVarP(&name, "name", "n", name, "Project name")
	newModelCmd.Flags().StringVarP(&projType, "projectType", "p", projType, "free, pro")
	newModelCmd.Flags().StringVarP(&modelType, "modelType", "m", modelType, "CAD, PHOTOGRAMMETRY, PTCLOUD")
	newModelCmd.Flags().StringVarP(&visibility, "visibility", "v", visibility, "public, unlisted, private")
	errors.Must(newModelCmd.MarkFlagRequired("name"))
}
