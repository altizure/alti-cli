package cmd

import (
	"fmt"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/spf13/cobra"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Username of current login user",
	Long:  `Show the username of the login user if user is loginned.`,
	Run: func(cmd *cobra.Command, args []string) {
		endpoint, user, err := gql.MySelf()
		if msg := errors.MustGQL(err, endpoint); msg != "" {
			fmt.Println(msg)
			return
		}
		if user.Username == "" {
			fmt.Printf("You are not login in!\nLogin with 'alti-cli login'\n")
		} else {
			fmt.Printf("%s: %s\n", endpoint, user.Username)
		}
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
