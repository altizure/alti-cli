package cmd

import (
	"fmt"
	"syscall"

	"github.com/jackytck/alti-cli/gql"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with email and password.",
	Long: `Login to Altizure with email and password.
	Credentials are stored in '~/.altizure/credentials'.`,
	Run: func(cmd *cobra.Command, args []string) {
		var email string
		fmt.Printf("Your login email: ")
		fmt.Scanln(&email)

		fmt.Printf("Your password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		password := string(bytePassword)
		fmt.Println()

		token := gql.GetUserToken(email, password, false)
		fmt.Printf("Token: %s\n", token)

		// @TODO: store token with end-point in ~/.altizure/credentials
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
