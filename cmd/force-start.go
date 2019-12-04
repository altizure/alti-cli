package cmd

import (
	"fmt"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/supergql"
	"github.com/spf13/cobra"
)

// forceStartCmd represents the forceStart command
var forceStartCmd = &cobra.Command{
	Use:   "force-start",
	Short: "Force start a project",
	Long:  "Force start a project.",
	Run: func(cmd *cobra.Command, args []string) {
		if !IsSuperuser() {
			fmt.Println("Not authorized.")
			return
		}

		p, err := gql.SearchProjectID(id, false)
		// input error
		if err != nil {
			fmt.Println("Project could not be found! Error:", err)
			return
		}

		tid, state, err := supergql.ForceStartTask(p.ID, taskType)
		// gql + sync error
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Successfully force start with tid: %q and state: %q\n", tid, state)
	},
}

func init() {
	superCmd.AddCommand(forceStartCmd)
	forceStartCmd.Flags().StringVarP(&id, "id", "p", id, "(Partial) Project id")
	forceStartCmd.Flags().StringVarP(&taskType, "type", "t", taskType, "Task type, default: Native")
	errors.Must(forceStartCmd.MarkFlagRequired("id"))
}
