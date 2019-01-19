package cmd

import (
	"fmt"
	"net/url"
	"syscall"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with email and password.",
	Long: `Login to Altizure with email and password.
	Credentials are stored in '~/.altizure/config'.`,
	Run: func(cmd *cobra.Command, args []string) {
		// a. api endpoint
		dc := config.DefaultConfig()
		dap := dc.GetActive()
		conf := config.Load()
		var endpoint string
		fmt.Printf("Endpoint (%s): ", dap.Endpoint)
		fmt.Scanln(&endpoint)
		if endpoint == "" {
			endpoint = dap.Endpoint
		}
		u, err := url.ParseRequestURI(endpoint)
		if err != nil {
			panic(err)
		}

		// b. api key
		var appKey string
		if u.Hostname() != config.DefaultHostName1 && u.Hostname() != config.DefaultHostName2 {
			fmt.Printf("App Key: ")
			fmt.Scanln(&appKey)
		} else {
			appKey = dap.Key
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

		token := gql.GetUserToken(endpoint, appKey, email, password, false)
		if token == "" {
			fmt.Println("Incorrect email or password!")
			return
		}

		// store key + token
		p := config.APoint{
			Endpoint: endpoint,
			Key:      appKey,
			Token:    token,
		}
		conf.AddProfile(p)
		err = conf.Save()
		errors.Must(err)

		// get username
		endpoint, user, err := gql.MySelf()
		errors.Must(err)

		// store username
		err = conf.SetActiveName(user.Username, true)
		errors.Must(err)

		fmt.Printf("Welcome %s (%s), you are logined to %s!\n", user.Name, user.Email, endpoint)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
