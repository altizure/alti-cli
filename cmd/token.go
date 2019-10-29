package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var email string

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Get user token by email",
	Long:  "Get user token by email.",
	Run: func(cmd *cobra.Command, args []string) {
		if !IsSuperuser() {
			fmt.Println("Not authorized.")
			return
		}
		config := config.Load()
		active := config.GetActive()
		token, err := gql.GetUserTokenBySuper(active.Endpoint, active.Key, active.Token, email)
		errors.Must(err)

		if jsonOut {
			o := tokenOutput{email, token}
			j, err := json.Marshal(o)
			errors.Must(err)
			js, err := gql.PrettyPrint(j)
			errors.Must(err)
			fmt.Println(js)
			return
		}

		if tableOut {
			header := []string{
				"Email",
				"Token",
			}
			row := []string{
				email,
				token,
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(header)
			table.Append(row)
			table.Render()
			return
		}

		fmt.Println(token)
	},
}

type tokenOutput struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func init() {
	superCmd.AddCommand(tokenCmd)
	tokenCmd.Flags().StringVarP(&email, "email", "e", email, "Email")
	tokenCmd.Flags().BoolVarP(&jsonOut, "json", "j", jsonOut, "Get JSON output.")
	tokenCmd.Flags().BoolVarP(&tableOut, "table", "t", tableOut, "Get table output.")
	errors.Must(tokenCmd.MarkFlagRequired("email"))
}
