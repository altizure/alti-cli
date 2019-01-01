package cmd

import (
	"fmt"
	"net/url"
	"syscall"

	"github.com/jackytck/alti-cli/config"
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
		// a. api endpoint
		var endpoint string
		fmt.Printf("Endpoint (%s): ", defaultEndpoint)
		fmt.Scanln(&endpoint)
		if endpoint == "" {
			endpoint = defaultEndpoint
		}
		u, err := url.ParseRequestURI(endpoint)
		if err != nil {
			panic(err)
		}

		// b. api key
		var appKey string
		if u.Hostname() != defaultHostName1 && u.Hostname() != defaultHostName2 {
			fmt.Printf("App Key: ")
			fmt.Scanln(&appKey)
		} else {
			appKey = defaultAppKey
		}

		// c. email
		var email string
		fmt.Printf("Your login email: ")
		fmt.Scanln(&email)

		// d. password
		fmt.Printf("Your password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		password := string(bytePassword)
		fmt.Println()

		token := gql.GetUserToken(endpoint, email, password, false)
		if token == "" {
			fmt.Println("Incorrect email or password!")
			return
		}

		c := config.Config{
			Endpoint: endpoint,
			Key:      appKey,
			Token:    token,
		}
		err = c.Save()
		if err != nil {
			panic(err)
		}
		fmt.Println("You are logined!")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
