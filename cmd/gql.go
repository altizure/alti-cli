package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/spf13/cobra"
)

var query string

// gqlCmd represents the gql command
var gqlCmd = &cobra.Command{
	Use:   "gql",
	Short: "Run arbitrary gql request.",
	Long:  "Run arbitrary gql request.",
	Run: func(cmd *cobra.Command, args []string) {
		q, err := ioutil.ReadFile(query)
		errors.Must(err)

		res, err := gql.Arbitrary(string(q))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(gqlCmd)
	gqlCmd.Flags().StringVarP(&query, "query", "q", query, "File storing the gql string.")
}
