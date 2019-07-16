package cmd

import (
	"fmt"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Choose which account to use.",
	Long: `List all accounts by 'account'. Get the reference ID and use this
command to switch to that, e.g. 'alti-cli account use ID'`,
	Run: func(cmd *cobra.Command, args []string) {
		config := config.Load()
		if len(args) < 1 {
			fmt.Println("Usage: alti-cli account use ID")
			return
		}
		id := args[0]
		p, err := config.GetProfile(id)
		if err != nil {
			fmt.Printf("ID: '%s' not found!\nLook up at 'alti-cli account'.", id)
			return
		}
		config.Active = p.ID
		err = config.Save()
		errors.Must(err)
		a := config.GetActive()
		fmt.Printf("Using: %s: %s\n", a.Endpoint, a.Name)
	},
}

func init() {
	accountCmd.AddCommand(useCmd)
}
