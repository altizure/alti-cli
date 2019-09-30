package cmd

import (
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/gql"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2/bson"
)

var actTimeout int

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
		timeout := time.Second * time.Duration(actTimeout)
		var accounts [][]string
		config := config.Load()
		for _, v := range config.Scopes {
			for _, p := range v.Profiles {
				var mode = gql.CheckSystemModeWithTimeout(v.Endpoint, p.Key, timeout)
				var info gql.AccountInfo
				var err error

				if mode == "Normal" {
					info, err = gql.GetAccountInfoTimeout(v.Endpoint, p.Key, p.Token, timeout)
				}
				if err != nil {
					mode = "Timeout"
				}

				nameOrEmail := p.Name
				if bson.IsObjectIdHex(p.Name) {
					nameOrEmail = p.Email
				}
				r := []string{p.ID, v.Endpoint, nameOrEmail, mode, "", "", "", "", "", "", "", ""}
				if config.Active == p.ID {
					r[4] = "Active"
				}
				if info.Sales {
					r[5] = "Yes"
				}
				if info.Super {
					r[6] = "Yes"
				}
				r[7] = strings.Join(info.ImageCloud, ",")
				r[8] = strings.Join(info.ModelCloud, ",")
				r[9] = strings.Join(info.MetaCloud, ",")
				r[10] = info.Version
				r[11] = info.ResponseTime.String()
				accounts = append(accounts, r)
			}
		}
		sort.Sort(byEndpoint(accounts))

		// render
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Endpoint", "Username/Email", "Status", "Select", "Sales", "Super", "Image Cloud", "Model Cloud", "Meta Cloud", "Version", "Response Time"})
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
	accountCmd.Flags().IntVarP(&actTimeout, "timeout", "t", 3, "Timeout of checking api server state in seconds")
}
