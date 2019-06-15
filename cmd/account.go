package cmd

import (
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/gql"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2/bson"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "List all the available accounts",
	Long:  "List all the previously logined accoutns across different servers.",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			elapsed := time.Since(start)
			log.Println("Took", elapsed)
		}()

		// prepare account list
		var accounts [][]string
		config := config.Load()
		for _, v := range config.Scopes {
			for _, p := range v.Profiles {
				var wg sync.WaitGroup
				wg.Add(5)

				var mode string
				var super, sales bool
				var cloud []string
				var version string

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
					cloud = gql.SupportedCloud(v.Endpoint, p.Key, "image")
					wg.Done()
				}()
				go func() {
					version = gql.Version(v.Endpoint, p.Key)
					wg.Done()
				}()

				wg.Wait()

				nameOrEmail := p.Name
				if bson.IsObjectIdHex(p.Name) {
					nameOrEmail = p.Email
				}
				r := []string{p.ID, v.Endpoint, nameOrEmail, mode, "", "", "", "", ""}
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
				r[8] = version
				accounts = append(accounts, r)
			}
		}
		sort.Sort(byEndpoint(accounts))

		// render
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Endpoint", "Username/Email", "Status", "Select", "Sales", "Super", "Image Upload Cloud", "Version"})
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
