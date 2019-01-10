package cmd

import (
	"fmt"
	"math"
	"os"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/types"
	tb "github.com/nsf/termbox-go"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// moreCmd represents the more command
var moreCmd = &cobra.Command{
	Use:   "more",
	Short: "List of all of my projects paginated.",
	Long:  "Show all of my projects in pages. Go to the next page by pressing n or Space or Enter. Previous page by p. Exit by q or Esc.",
	Run: func(cmd *cobra.Command, args []string) {
		err := tb.Init()
		errors.Must(err)
		fmt.Println("Loading...")
		page, total, table, err := get(pageCount, 0, "", "")
		// so that when termbox is quit, the last rendered result could be shown
		defer func() {
			tb.Close()
			table.Render()
		}()
		errors.Must(err)
		clear()
		fmt.Printf("Totals: %d (Next: n or Space or Enter. Previous: p. Exit: q or Esc.)\n", total)
		table.Render()
		if !page.HasNextPage {
			return
		}
		curPage := 0
		maxPage := int(math.Ceil(float64(total) / float64(pageCount)))
		fmt.Printf("Page: %d/%d\n", curPage+1, maxPage)
		for {
			evt := tb.PollEvent()
			switch {
			case curPage+1 < maxPage && (evt.Ch == 'n' || evt.Key == tb.KeySpace || evt.Key == tb.KeyEnter):
				page, _, table, err = next(page.EndCursor)
				errors.Must(err)
				clear()
				table.Render()
				curPage++
				fmt.Printf("Page: %d/%d\n", curPage+1, maxPage)
			case curPage > 0 && evt.Ch == 'p':
				page, _, table, err = prev(page.StartCursor)
				errors.Must(err)
				clear()
				table.Render()
				curPage--
				fmt.Printf("Page: %d/%d\n", curPage+1, maxPage)
			case evt.Ch == 'q' || evt.Key == tb.KeyEsc || evt.Key == tb.KeyCtrlC:
				return
			}
		}
	},
}

func clear() {
	tb.Clear(tb.ColorWhite, tb.ColorBlack)
	tb.Flush()
}

func next(cur string) (*types.PageInfo, int, *tablewriter.Table, error) {
	return get(pageCount, 0, "", cur)
}

func prev(cur string) (*types.PageInfo, int, *tablewriter.Table, error) {
	return get(0, pageCount, cur, "")
}

func get(first, last int, before, after string) (*types.PageInfo, int, *tablewriter.Table, error) {
	projs, page, total, err := gql.MyProjects(first, last, before, after, search)
	if err != nil {
		return nil, 0, nil, err
	}
	table := types.ProjectsToTable(projs, os.Stdout)
	return page, total, table, nil
}

func init() {
	myprojCmd.AddCommand(moreCmd)
	moreCmd.Flags().IntVarP(&pageCount, "count", "c", pageCount, "number of projects per page")
	moreCmd.Flags().StringVarP(&search, "search", "q", search, "display name to search")
}
