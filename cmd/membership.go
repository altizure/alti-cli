package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var jsonOut bool

// membershipCmd represents the membership command
var membershipCmd = &cobra.Command{
	Use:   "membership",
	Short: "Membership info",
	Long:  "My membership info if any.",
	Run: func(cmd *cobra.Command, args []string) {
		endpoint, user, err := gql.MySelf()
		if msg := errors.MustGQL(err, endpoint); msg != "" {
			fmt.Println(msg)
			return
		}
		if user.Username == "" {
			fmt.Printf("You are not login in!\nLogin with 'alti-cli login'\n")
			return
		}

		m := user.Membership
		if jsonOut {
			j, err := json.Marshal(m)
			errors.Must(err)
			js, err := gql.PrettyPrint(j)
			errors.Must(err)
			fmt.Println(js)
			return
		}

		header := []string{
			"State",
			"Plan",
			"Months",
			"Start",
			"End",
			"GP Quota",
			"Coin/GP",
			"Storage",
			"Visibility",
			"Coupon",
			"Model/Project",
			"Collaborator",
			"Watermark",
		}
		row := []string{
			m.State,
			m.PlanName,
			strconv.Itoa(m.Period),
			m.StartDate.Format("2006-01-02 15:04:05"),
			m.EndDate.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%.2f", m.MemberGPQuota),
			fmt.Sprintf("%.2f", m.CoinPerGP),
			fmt.Sprintf("%.2f", m.AssetStorage),
			strings.Join(m.Visibility, ", "),
			m.Coupon.String(),
			strconv.Itoa(m.ModelPerProject),
			strconv.Itoa(m.CollaboratorQuota),
			fmt.Sprintf("%v", m.ForceWatermark),
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(header)
		table.Append(row)
		table.Render()
	},
}

func init() {
	myCmd.AddCommand(membershipCmd)
	membershipCmd.Flags().BoolVarP(&jsonOut, "json", "j", jsonOut, "Get JSON output.")
}
