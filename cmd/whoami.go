package cmd

import (
	"fmt"

	"github.com/jackytck/alti-cli/gql"
	"github.com/spf13/cobra"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Username of current login user",
	Long:  `Show the username of the login user if user is loginned.`,
	Run: func(cmd *cobra.Command, args []string) {
		user, err := gql.MySelf()
		if err != nil {
			panic(err)
		}
		if user.Username == "" {
			fmt.Printf("You are not login in!\nLogin with 'alti-cli login'\n")
		} else {
			fmt.Println(user.Username)
		}
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
