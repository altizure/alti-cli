package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/spf13/cobra"
)

var queryFile string
var varFile string

// gqlCmd represents the gql command
var gqlCmd = &cobra.Command{
	Use:   "gql",
	Short: "Run arbitrary gql request.",
	Long:  "Run arbitrary gql request.",
	Run: func(cmd *cobra.Command, args []string) {
		q, err := ioutil.ReadFile(queryFile)
		if err != nil {
			fmt.Println(errors.ErrClientQuery)
			return
		}
		va := make(map[string]interface{})
		if varFile != "" {
			vb, err2 := ioutil.ReadFile(varFile)
			if err2 != nil {
				fmt.Println(errors.ErrClientVar)
				return
			}
			err = json.Unmarshal(vb, &va)
			if err != nil {
				fmt.Println(errors.ErrClientVarInvalid)
				return
			}
		}

		res, err := gql.Arbitrary(string(q), va)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(gqlCmd)
	gqlCmd.Flags().StringVarP(&queryFile, "query", "q", queryFile, "File storing the gql string.")
	gqlCmd.Flags().StringVarP(&varFile, "variable", "k", varFile, "File storing the related variables.")
}
