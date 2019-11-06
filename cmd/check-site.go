package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackytck/alti-cli/cloud"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/spf13/cobra"
)

// checkSiteCmd represents the checkSite command
var checkSiteCmd = &cobra.Command{
	Use:   "site",
	Short: "Check main site availability",
	Long:  "Check main site availability",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			if verbose {
				elapsed := time.Since(start)
				log.Println("Took", elapsed)
			}
		}()

		if err := service.Check(
			nil,
			service.CheckAPIServer(),
		); err != nil {
			log.Println(err)
			return
		}

		url := gql.WebEndpoint()
		if verbose {
			log.Printf("Checking %s...\n", url)
		}

		ch := make(chan res)
		go func() {
			s, err := cloud.GetStatus(url)
			ch <- res{s, err}
		}()

		select {
		case <-time.After(time.Second * time.Duration(timeout)):
			log.Println(errors.ErrClientTimeout)
		case r := <-ch:
			if r.err != nil {
				log.Println(r.err)
				return
			}
			if r.status != http.StatusOK {
				if verbose {
					log.Printf("Status is not OK: %v", r.status)
				}
				return
			}
			if verbose {
				log.Printf("Success with status code: %v", r.status)
			} else {
				fmt.Println("Success")
			}
		}
	},
}

type res struct {
	status int
	err    error
}

func init() {
	checkCmd.AddCommand(checkSiteCmd)
	checkSiteCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "Timeout of checking in seconds")
	checkSiteCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display more info of operation")
}
