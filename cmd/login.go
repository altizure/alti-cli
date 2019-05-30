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

var byPhone bool

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

		var token string

		if byPhone {
			// c1. phone
			var phone string
			fmt.Printf("Your phone number (international): ")
			fmt.Scanln(&phone)

			err = gql.RequestLoginCode(endpoint, appKey, phone)
			if err != nil {
				panic(err)
			}

			// d1. login code from sms
			var code string
			fmt.Printf("Your login code: ")
			fmt.Scanln(&code)

			token, err = gql.GetUserTokenByCode(endpoint, appKey, phone, code)
			if err != nil || token == "" {
				fmt.Println("Incorrect code!")
				return
			}
		} else {
			// c2. email
			var email string
			fmt.Printf("Your login email: ")
			fmt.Scanln(&email)

			// d2. password
			fmt.Printf("Your password: ")
			bytePassword, err2 := terminal.ReadPassword(int(syscall.Stdin))
			if err2 != nil {
				panic(err2)
			}
			password := string(bytePassword)
			fmt.Println()

			token, err = gql.GetUserTokenByEmail(endpoint, appKey, email, password, false)
			if err != nil || token == "" {
				fmt.Println("Incorrect email or password!")
				return
			}
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
		err = conf.SetActiveUserInfo(user.Username, user.Email, true)
		errors.Must(err)

		fmt.Printf("Welcome %s (%s), you are logined to %s!\n", user.Name, user.Email, endpoint)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().BoolVarP(&byPhone, "phone", "p", byPhone, "Use verified phone number to login")
}
