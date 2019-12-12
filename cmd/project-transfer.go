package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/spf13/cobra"
)

// projTransferCmd represents the project transfer command
var projTransferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer my project to another user.",
	Long:  "Transfer my project to another user by email, with custom message.",
	Run: func(cmd *cobra.Command, args []string) {
		// pre-checks general
		if err := service.Check(
			nil,
			service.CheckAPIServer(),
			service.CheckPID("", id),
		); err != nil {
			log.Println(err)
			return
		}

		// project info
		p, err := gql.Project(id)
		if err != nil {
			log.Println(err)
			return
		}

		// confirm?
		var ans string
		fmt.Printf("Are you sure to transfer project: %q (%s) to  %q? (Y/N): ", p.Name, id, email)
		if assumeYes {
			fmt.Println("Yes")
		} else {
			fmt.Scanln(&ans)
			ans = strings.ToUpper(ans)
			if ans != "Y" && ans != service.Yes {
				log.Println("Cancelled.")
				return
			}
		}

		res, err := gql.TransferProject(id, email, message)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Successfully transferred project: %q (%s) to %q with status: %q\n", p.Name, id, email, res)
	},
}

func init() {
	projectCmd.AddCommand(projTransferCmd)
	projTransferCmd.Flags().StringVarP(&id, "id", "p", id, "(Partial) Project id")
	projTransferCmd.Flags().StringVarP(&email, "email", "e", email, "Recipient email")
	projTransferCmd.Flags().StringVarP(&message, "message", "m", message, "Message to recipient")
	projTransferCmd.Flags().BoolVarP(&assumeYes, "assumeyes", "y", assumeYes, "Assume yes; assume that the answer to any question which would be asked is yes")
	errors.Must(projTransferCmd.MarkFlagRequired("id"))
	errors.Must(projTransferCmd.MarkFlagRequired("email"))
}
