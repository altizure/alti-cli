package cmd

import (
	"os"
	"strconv"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/web"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Check if api server could reach this client",
	Long:  "Locally start a web server and check if the api server could reach this server.",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := web.CheckVisibility()
		errors.Must(err)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"URL", "Visibility"})
		for k, v := range res {
			r := []string{k, strconv.FormatBool(v)}
			table.Append(r)
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)
}
