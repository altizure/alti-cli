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

var tableOut bool

// myTokenCmd represents the my token command
var myTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Get my user token",
	Long:  "Get my currently active user token with respect to the app.",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.Load()
		active := config.GetActive()

		if active.Token == "" {
			fmt.Println("Please login with:\nalti-cli login")
			return
		}

		if jsonOut {
			j, err := json.Marshal(active)
			errors.Must(err)
			js, err := gql.PrettyPrint(j)
			errors.Must(err)
			fmt.Println(js)
			return
		}

		if tableOut {
			header := []string{
				"Endpoint",
				"Name",
				"Email",
				"Key",
				"Token",
			}
			row := []string{
				active.Endpoint,
				active.Name,
				active.Email,
				active.Key,
				active.Token,
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(header)
			table.Append(row)
			table.Render()
			return
		}

		fmt.Println(active.Token)
	},
}

func init() {
	myCmd.AddCommand(myTokenCmd)
	myTokenCmd.Flags().BoolVarP(&jsonOut, "json", "j", jsonOut, "Get JSON output.")
	myTokenCmd.Flags().BoolVarP(&tableOut, "table", "t", tableOut, "Get table output.")
}
