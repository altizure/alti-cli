package cmd

import (
	"fmt"
	"log"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/spf13/cobra"
)

var desc string

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report issue of a project.",
	Run: func(cmd *cobra.Command, args []string) {
		// pre-checks
		if err := service.Check(
			nil,
			service.CheckAPIServer(),
			service.CheckPID("", id),
		); err != nil {
			log.Println(err)
			return
		}

		// get pid
		p, _ := gql.SearchProjectID(id, true)

		errors.Must(gql.ReportProject(p.ID, desc))

		fmt.Printf("Successfully reported project: %q", p.ID)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringVarP(&id, "id", "p", id, "Project (partial) id")
	reportCmd.Flags().StringVarP(&desc, "desc", "d", desc, "Description of the issue")
	errors.Must(reportCmd.MarkFlagRequired("id"))
	errors.Must(reportCmd.MarkFlagRequired("desc"))
}
