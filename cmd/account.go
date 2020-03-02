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
	"github.com/jackytck/alti-cli/service"
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
		cfg := config.Load()
		timeout := time.Second * time.Duration(actTimeout)
		actCh := make(chan []string)
		var wg sync.WaitGroup
		wg.Add(cfg.Size())

		go func() {
			defer func() {
				close(actCh)
			}()
			wg.Wait()
		}()

		for _, s := range cfg.Scopes {
			for _, p := range s.Profiles {
				go func(s config.Scope, p config.Profile) {
					actCh <- getActRow(s, p, cfg.Active, timeout)
					wg.Done()
				}(s, p)
			}
		}

		// aggregate all accounts
		var accounts [][]string
		for r := range actCh {
			accounts = append(accounts, r)
		}
		sort.Sort(byEndpoint(accounts))

		// render
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Endpoint", "Username/Email", "Status", "Select", "Sales", "Super", "Image Cloud", "Model Cloud", "Meta Cloud", "Version", "Response Time"})
		table.AppendBulk(accounts)
		table.Render()

		// readme
		log.Println("To switch account: Use: alti-cli account use ID")
	},
}

func getActRow(s config.Scope, p config.Profile, activeID string, timeout time.Duration) []string {
	var mode = gql.CheckSystemModeWithTimeout(s.Endpoint, p.Key, timeout)
	var info gql.AccountInfo
	var err error

	if mode == service.NormalMode || mode == service.ReadOnlyMode {
		info, err = gql.GetAccountInfoTimeout(s.Endpoint, p.Key, p.Token, timeout)
	}
	if err != nil {
		mode = "Timeout"
	}

	nameOrEmail := p.Name
	if p.Name == "" || bson.IsObjectIdHex(p.Name) {
		nameOrEmail = p.Email
	}

	r := []string{p.ID, s.Endpoint, nameOrEmail, mode, "", "", "", "", "", "", "", ""}
	if p.ID == activeID {
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

	return r
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
