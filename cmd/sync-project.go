package cmd

import (
	"fmt"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/supergql"
	"github.com/spf13/cobra"
)

// tokenCmd represents the token command
var syncProjCmd = &cobra.Command{
	Use:   "sync",
	Short: "Cloud sync a project",
	Long:  "Cloud sync a project.",
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

		state, err := supergql.SyncProject(p.ID)
		// gql + sync error
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Successfully triggered with state: %q\n", state)
	},
}

func init() {
	superCmd.AddCommand(syncProjCmd)
	syncProjCmd.Flags().StringVarP(&id, "id", "p", id, "(Partial) Project id")
	errors.Must(syncProjCmd.MarkFlagRequired("id"))
}
