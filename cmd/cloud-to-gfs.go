package cmd

import (
	"fmt"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/supergql"
	"github.com/spf13/cobra"
)

// cloudToGFSCmd represents the cloudToGFS command
var cloudToGFSCmd = &cobra.Command{
	Use:   "toGFS",
	Short: "Download image from cloud to gfs",
	Long:  "Download image from cloud to gfs.",
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

		state, err := supergql.TriggerCloudToGFS(p.ID)
		// gql + trigger error
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Successfully triggered with state: %q\n", state)
	},
}

func init() {
	superCmd.AddCommand(cloudToGFSCmd)
	cloudToGFSCmd.Flags().StringVarP(&id, "id", "p", id, "(Partial) Project id")
	errors.Must(cloudToGFSCmd.MarkFlagRequired("id"))
}
