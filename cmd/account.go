package cmd

import (
	"os"
	"sort"
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
		// prepare account list
		var accounts [][]string
		config := config.Load()
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
				accounts = append(accounts, r)
			}
		}
		sort.Sort(byEndpoint(accounts))

		// render
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Endpoint", "Username", "Status", "Select", "Sales", "Super", "Upload Cloud"})
		table.AppendBulk(accounts)
		table.Render()
	},
}

type byEndpoint [][]string

func (e byEndpoint) Len() int {
	return len(e)
}

func (e byEndpoint) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e byEndpoint) Less(i, j int) bool {
	return e[i][1]+e[i][0] < e[j][1]+e[j][0]
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
