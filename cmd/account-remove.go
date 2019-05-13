package cmd

import (
	"fmt"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove account profile",
	Long:  "Remove a non-default user profile from the account list and config file.",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.Load()
		if len(args) < 1 {
			fmt.Println("Usage: alti-cli account remove ID")
			return
		}
		id := args[0]

		p, err := config.RemoveProfile(id, true)
		if err != nil {
			switch err {
			case errors.ErrProfileNotFound:
				fmt.Printf("ID: '%s' not found!\nLook up at 'alti-cli account'.", id)
			case errors.ErrProfileNotRemovable:
				fmt.Println("Default profile could not be removed!")
			}
			return
		}

		fmt.Printf("Removed: %s: %s\n", p.Endpoint, p.Name)
	},
}

func init() {
	accountCmd.AddCommand(removeCmd)
}
