package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackytck/alti-cli/cloud"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/file"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var projDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download reconstruction results if any.",
	Long:  "Download any project reconstruction resutls into current directory.",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			if verbose {
				elapsed := time.Since(start)
				log.Println("Took", elapsed)
			}
		}()

		p, err := gql.SearchProjectID(id, false)
		if err != nil {
			fmt.Println("Project could not be found! Error:", err)
			return
		}

		// display
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"PID", "State", "Name", "Size", "Last modified", "Link"})
		var items []item
		for _, d := range p.Downloads.Edges {
			state := d.Node.State
			name := d.Node.Name
			size := fmt.Sprintf("%.2f MB", file.BytesToMB(d.Node.Size))
			modified := d.Node.Mtime.Format("2006-01-02 15:04:05")
			link := d.Node.Link
			table.Append([]string{p.ID, state, name, size, modified, link})
			if link != "" {
				items = append(items, item{name, link})
			}
		}
		table.Render()

		// download
		total := len(items)
		if total == 0 {
			fmt.Println("No downloadable could be found!")
			return
		}

		// ask user to proceed or not
		var ans string
		plural := ""
		if total > 1 {
			plural = "s"
		}
		fmt.Printf("Continue to download %d item%s or not? (Y/N): ", total, plural)
		if assumeYes {
			fmt.Println("Yes")
		} else {
			fmt.Scanln(&ans)
			ans = strings.ToUpper(ans)
			if ans != "Y" && ans != service.Yes {
				log.Println("Cancelled.")
				return
			}
		}

		for _, v := range items {
			log.Printf("Downloading %q...", v.Name)
			errors.Must(cloud.GetFile(v.Name, v.Link))
			log.Println("Done")
		}
	},
}

type item struct {
	Name string
	Link string
}

func init() {
	projectCmd.AddCommand(projDownloadCmd)
	projDownloadCmd.Flags().StringVarP(&id, "id", "p", id, "(Partial) Project id")
	projDownloadCmd.Flags().BoolVarP(&assumeYes, "assumeyes", "y", assumeYes, "Assume yes; assume that the answer to any question which would be asked is yes")
	projDownloadCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display more info")
	errors.Must(projDownloadCmd.MarkFlagRequired("id"))
}
