package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// bucketCmd represents the bucket command
var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "List all available buckets",
	Long:  `'alti-cli list bucket' to list all available buckets of different types.`,
	Run: func(cmd *cobra.Command, args []string) {
		clouds := gql.SupportedCloud("", "")
		kinds := []string{"image", "model"}

		var buckets [][]string
		for _, k := range kinds {
			for _, c := range clouds {
				buks, err := gql.BucketList(k, c)
				if err != nil {
					if err != errors.ErrBucketInvalid {
						panic(err)
					}
					continue
				}
				buckets = append(buckets, []string{k, strings.ToLower(c), strings.Join(buks, ", "), fmt.Sprintf("%d", len(buks))})
			}
		}

		// render
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Kind", "Cloud", "Buckets", "Count"})
		table.AppendBulk(buckets)
		table.Render()
	},
}

func init() {
	listCmd.AddCommand(bucketCmd)
}
