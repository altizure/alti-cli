package cmd

import (
	"os"
	"strings"
	"sync"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/gql"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "List all the available accounts",
	Long:  "List all the previously logined accoutns across different servers.",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.Load()
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Endpoint", "Username", "Status", "Select", "Sales", "Super", "Upload Cloud"})
		for _, v := range config.Scopes {
			for _, p := range v.Profiles {
				var wg sync.WaitGroup
				wg.Add(4)

				var mode string
				var super, sales bool
				var cloud []string

				go func() {
					mode = gql.CheckSystemMode(v.Endpoint, p.Key)
					wg.Done()
				}()
				go func() {
					sales = gql.IsSales(v.Endpoint, p.Key, p.Token)
					wg.Done()
				}()
				go func() {
					super = gql.IsSuper(v.Endpoint, p.Key, p.Token)
					wg.Done()
				}()
				go func() {
					cloud = gql.SupportedCloud(v.Endpoint, p.Key)
					wg.Done()
				}()

				wg.Wait()
				r := []string{p.ID, v.Endpoint, p.Name, mode, "", "", "", ""}
				if config.Active == p.ID {
					r[4] = "Active"
				}
				if sales {
					r[5] = "Yes"
				}
				if super {
					r[6] = "Yes"
				}
				r[7] = strings.Join(cloud, ",")
				table.Append(r)
			}
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
