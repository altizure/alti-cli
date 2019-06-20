package cmd

import (
	"log"
	"os"
	"time"

	"github.com/jackytck/alti-cli/gql"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var errCode string
var errLang = "en"

// errorCmd represents the error command
var errorCmd = &cobra.Command{
	Use:   "error",
	Short: "Query description and solution of error code.",
	Long:  "Query the description and suggested solution of an error code.",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			if verbose {
				elapsed := time.Since(start)
				log.Println("Took", elapsed)
			}
		}()

		info, err := gql.GetErrorCodeInfo(errCode, errLang)
		if err != nil {
			log.Printf("Invalid code: %q\n", errCode)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Append([]string{"Code", info.Code})
		table.Append([]string{"Description", info.Description})
		table.Append([]string{"Solution", info.Solution})
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(errorCmd)
	errorCmd.Flags().StringVarP(&errCode, "code", "c", errCode, "Altizure error code.")
	errorCmd.Flags().StringVarP(&errLang, "lang", "l", errLang, "Language of error info. One of 'en', 'zh_tw' and 'zh_cn'.")
	errorCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display more info of operation")
}
