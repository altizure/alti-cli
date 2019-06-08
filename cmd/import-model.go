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

	"github.com/jackytck/alti-cli/cloud"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/jackytck/alti-cli/web"
	"github.com/spf13/cobra"
)

var model string
var ip string
var port string
var timeout int
var partsDir string

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
			service.CheckZip(model),
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

		// register + upload + state check
		mru := cloud.ModelRegUploader{
			Method:       method,
			PID:          proj.ID,
			ModelPath:    model,
			Filename:     filename,
			DirectURL:    directURL,
			Bucket:       bucket,
			MultipartDir: partsDir,
			Verbose:      verbose,
		}
		state, err := mru.Run()
		if err != nil {
			log.Printf(err.Error())
			return
		}

		log.Printf("Successfully registered and uplaoded in state: %q!\n", state)
	},
}

func init() {
	importCmd.AddCommand(importModelCmd)
	importModelCmd.Flags().StringVarP(&id, "id", "p", id, "Project id")
	importModelCmd.Flags().StringVarP(&model, "file", "f", model, "File path of model zip file.")
	importModelCmd.Flags().StringVarP(&method, "method", "m", method, "Desired method of upload: 'direct' or 's3'")
	importModelCmd.Flags().IntVarP(&timeout, "timeout", "t", timeout, "Timeout of checking direct upload state in seconds")
	importModelCmd.Flags().StringVar(&ip, "ip", ip, "IP address of ad-hoc local server for direct upload.")
	importModelCmd.Flags().StringVar(&port, "port", port, "Port of ad-hoc local server for direct upload.")
	importModelCmd.Flags().StringVarP(&bucket, "bucket", "b", bucket, "Desired bucket to upload for method: 's3'")
	importModelCmd.Flags().StringVarP(&partsDir, "parts", "d", partsDir, "Directory of 7z splitted multiparts.")
	importModelCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display more info of operation")
	importModelCmd.MarkFlagRequired("id")
	importModelCmd.MarkFlagRequired("file")
}
