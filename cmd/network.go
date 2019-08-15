package cmd

import (
	"fmt"
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
		u, res, err := web.PreferredLocalURL(verbose)
		if err != errors.ErrClientInvisible {
			errors.Must(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"URL", "Visibility"})
		for k, v := range res {
			r := []string{k, strconv.FormatBool(v)}
			table.Append(r)
		}
		table.Render()

		if err == errors.ErrClientInvisible {
			fmt.Println("Client is invisible. Direct upload is not supported!")
			return
		}
		fmt.Printf("Preferred %q for direct upload!\n", u.Hostname())
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)
	networkCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display more network checking info")
}
