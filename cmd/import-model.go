package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/file"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/jackytck/alti-cli/web"
	"github.com/spf13/cobra"
)

var model string
var ip string
var port string

// importModelCmd represents the importModel command
var importModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Import model from a local / remote path into a project",
	Long:  "Check and upload third party model into a project.",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			if verbose {
				elapsed := time.Since(start)
				log.Println("Took", elapsed)
			}
		}()

		// pre-checks
		if err := service.Check(
			nil,
			service.CheckAPIServer(),
			service.CheckUploadMethod("model", method, ip, port),
			service.CheckPID("model", id),
			service.CheckFile(model),
		); err != nil {
			log.Println(err)
			return
		}

		// get project
		proj, _ := gql.SearchProjectID(id, true)

		// setup direct upload server
		method = strings.ToLower(method)
		var serDone func()
		var baseURL, directURL string
		filename := filepath.Base(model)
		if method == service.DirectUploadMethod {
			bu, done, err := web.StartLocalServer(filepath.Dir(model), ip, port, false)
			errors.Must(err)
			defer done()
			serDone = done
			baseURL = bu
			directURL = fmt.Sprintf("%s/%s", baseURL, filename)
		}

		// capture ctrl+c
		cc := make(chan os.Signal)
		signal.Notify(cc, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-cc
			fmt.Println()
			serDone()
			log.Println("Bye!")
			os.Exit(1)
		}()

		// compute checksum
		log.Printf("Computing checksum: %q...\n", filename)
		checksum, err := file.Sha1sum(model)
		errors.Must(err)
		log.Printf("SHA1: %s\n", checksum)

		// register model
		if method == service.DirectUploadMethod {
			im, err := gql.RegisterModelURL(proj.ID, directURL, filename, checksum)
			errors.Must(err)
			log.Printf("Registered model with state: %q\n", im.State)

			// check if ready
			stateC := make(chan string)

			go func() {
				log.Println("Checking state...")
				for {
					p, err := gql.Project(proj.ID)
					errors.Must(err)
					s := p.ImportedState
					if s != "Pending" {
						stateC <- s
						return
					}
					time.Sleep(time.Second * 1)
				}
			}()

			var state string
			select {
			case <-time.After(time.Minute * 1):
				log.Printf("Client timeout.")
				return
			case state = <-stateC:
				log.Printf("Model is in state: %q\n", state)
			}
		}
	},
}

func init() {
	importCmd.AddCommand(importModelCmd)
	importModelCmd.Flags().StringVarP(&id, "id", "p", id, "Project id")
	importModelCmd.Flags().StringVarP(&model, "file", "f", model, "File path of model zip file.")
	importModelCmd.Flags().StringVarP(&method, "method", "m", method, "Desired method of upload: 'direct' or 's3'")
	importModelCmd.Flags().StringVar(&ip, "ip", ip, "IP address of ad-hoc local server for direct upload.")
	importModelCmd.Flags().StringVar(&port, "port", port, "Port of ad-hoc local server for direct upload.")
	importModelCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display more info of operation")
	importModelCmd.MarkFlagRequired("id")
	importModelCmd.MarkFlagRequired("file")
}
