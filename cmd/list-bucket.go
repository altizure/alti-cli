package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// bucketCmd represents the bucket command
var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "List all available buckets",
	Long:  `'alti-cli list bucket' to list all available buckets of different types.`,
	Run: func(cmd *cobra.Command, args []string) {
		// pre-checks
		if err := service.Check(
			nil,
			service.CheckAPIServer(),
		); err != nil {
			log.Println(err)
			return
		}

		kinds := []string{"image", "meta", "model"}

		var buckets [][]string
		for _, k := range kinds {
			clouds := gql.SupportedCloud("", "", k)
			for _, c := range clouds {
				buks, err := gql.BucketList(k, c)
				if err != nil {
					if err != errors.ErrBucketInvalid {
						panic(err)
					}
					continue
				}
				suggested, err := gql.SuggestedBucket(k, c)
				if err != nil {
					panic(err)
				}
				buckets = append(buckets, []string{k, strings.ToLower(c), strings.Join(buks, ", "), suggested, fmt.Sprintf("%d", len(buks))})
			}
		}

		// render
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Kind", "Cloud", "Buckets", "Suggested", "Count"})
		table.AppendBulk(buckets)
		table.Render()
	},
}

func init() {
	listCmd.AddCommand(bucketCmd)
}
