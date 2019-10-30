package cmd

import (
	"fmt"
	"os"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/types"
	"github.com/olekukonko/tablewriter"
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
			fmt.Println(LoginHint)
			return
		}

		header := []string{"Endpoint"}
		row := []string{endpoint}
		header = append(header, types.UserHeaderString()...)
		row = append(row, user.RowString()...)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(header)
		table.Append(row)
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
