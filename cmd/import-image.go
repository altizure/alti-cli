package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/jackytck/alti-cli/cloud"
	"github.com/jackytck/alti-cli/db"
	"github.com/jackytck/alti-cli/file"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/types"
	"github.com/jackytck/alti-cli/web"
	"github.com/spf13/cobra"
)

var method string
var report string

// importImageCmd represents the importImage command
var importImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Import images from a directory into a project",
	Long:  "Check and upload images into a project.",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			if verbose {
				elapsed := time.Since(start)
				log.Println("Took", elapsed)
			}
		}()

		// check pid
		p, err := gql.SearchProjectID(id, true)
		if err != nil {
			fmt.Println("Project could not be found! Error:", err)
			return
		}
		log.Printf("Importing to %q with pid: %q...\n", p.Name, p.ID)

		log.Printf("Checking %s...\n", dir)

		// stats
		var totalGP float64
		var totalImg int
		var totalByte datasize.ByteSize
		var existedCnt int

		// setup image digester
		done := make(chan struct{})
		defer close(done)

		paths, errc := file.WalkFiles(done, dir, skip)
		result := make(chan file.ImageDigest)

		digester := file.ImageDigester{
			Root:   dir,
			PID:    p.ID,
			Done:   done,
			Paths:  paths,
			Result: result,
		}
		threads := digester.Run(thread)
		if verbose {
			log.Printf("Working in %d thread(s)...", threads)
		}

		// setup local temp db
		dbPath, err := db.OpenPath()
		if err != nil {
			panic(err)
		}
		localDB, err := db.OpenDB(dbPath)
		if err != nil {
			panic(err)
		}
		err = localDB.Init(&db.Image{})
		if err != nil {
			panic(err)
		}
		cleanupDB := func() {
			err2 := os.Remove(dbPath)
			if err2 != nil {
				panic(err2)
			}
		}
		defer cleanupDB()

		// capture ctrl+c
		cc := make(chan os.Signal)
		signal.Notify(cc, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-cc
			cleanupDB()
			fmt.Println()
			log.Println("Bye!")
			os.Exit(1)
		}()

		for r := range result {
			if r.Error != nil {
				log.Printf("Invalid image: %q, Reason: %v", r.Path, r.Error)
				continue
			}

			mb := file.BytesToMB(r.Filesize)
			if verbose {
				log.Printf("Path: %q, URL: %q, Filename: %q, Dimension: %d x %d, GP: %.2f, Type: %s, Size: %.2f MB, Checksum: %s, Existed: %v\n",
					r.Path, r.URL, r.Filename, r.Width, r.Height, r.GP, r.Filetype, mb, r.SHA1, r.Existed)
			}

			if r.Existed {
				existedCnt++
				continue
			}

			totalGP += r.GP
			totalImg++
			totalByte += datasize.ByteSize(r.Filesize)

			img := db.Image{
				PID:       p.ID,
				Filename:  r.Filename,
				Filetype:  types.ConvertToImageType(r.Filetype),
				URL:       r.URL,
				LocalPath: r.Path,
				Hash:      r.SHA1,
				Width:     r.Width,
				Height:    r.Height,
				GP:        r.GP,
			}
			err = localDB.Save(&img)
			if err != nil {
				panic(err)
			}
		}

		// check whether the Walk failed
		if err = <-errc; err != nil {
			panic(err)
		}

		if totalImg == 0 {
			if existedCnt > 0 {
				log.Println("No new image is found! All of the images in this directory have been imported.")
			} else {
				log.Println("No image is found!")
			}
			return
		}

		// ask user to proceed or not
		var ans string
		if existedCnt > 0 {
			log.Printf("%d images already existed in the project", existedCnt)
		}
		log.Printf("Found %d images, total %.2f GP, %s", totalImg, totalGP, totalByte.HumanReadable())
		plural := ""
		if totalImg > 1 {
			plural = "s"
		}
		fmt.Printf("After importing (if no duplicate):\nImages #: %d -> %d\tGP: %.2f -> %.2f\n", p.NumImage, p.NumImage+totalImg, p.GigaPixel, p.GigaPixel+totalGP)
		fmt.Printf("Continue to import %d image%s or not? (Y/N): ", totalImg, plural)
		fmt.Scanln(&ans)
		ans = strings.ToUpper(ans)
		if ans != "Y" && ans != "YES" {
			log.Println("Cancelled.")
			return
		}

		// check direct upload
		var baseURL string
		if method == "" {
			method = "direct"
			var localSever *http.Server
			var port int
			log.Println("Checking direct upload...")
			pu, _, err2 := web.PreferedLocalURL(verbose)
			baseURL = ""
			if err2 != nil {
				log.Println("Client is invisible. Direct upload is not supported!")
				log.Println("Using S3.")
				method = "s3"
			} else {
				log.Printf("Direct upload is supported over %q!\n", pu.Hostname())

				// setup local web server
				s := web.Server{Directory: dir, Address: pu.Hostname() + ":"}
				fmt.Println("local server", s)
				localSever, port, err = s.ServeStatic(verbose)
				if err != nil {
					panic(err)
				}
				baseURL = fmt.Sprintf("http://%s:%d", pu.Hostname(), port)
				log.Printf("Serving files at %s\n", baseURL)

				defer func() {
					if err = localSever.Shutdown(context.TODO()); err != nil {
						panic(err)
					}
				}()
			}
		} else {
			log.Printf("Using %s to upload\n", method)
			bucket, err2 := gql.BucketSuggestion(method)
			if err2 != nil {
				panic(err2)
			}
			log.Printf("Bucket %q is chosen", bucket)
		}

		// read from local db, register and upload
		imgc, errc := db.AllImage(localDB)
		ruRes := make(chan db.Image)
		ruDigester := cloud.ImageRegUploader{
			Method:  method,
			BaseURL: baseURL,
			Images:  imgc,
			Done:    done,
			Result:  ruRes,
		}
		ruDigester.Run(thread)

		for img := range ruRes {
			err = localDB.Save(&img)
			if verbose {
				if img.Error != "" {
					log.Printf("Registration failed: %q\n", img.Error)
				} else {
					log.Printf("Registered %q\n", img.Filename)
				}
			}
			if err != nil {
				panic(err)
			}
		}

		// check whether the read from local db failed
		if err = <-errc; err != nil {
			panic(err)
		}

		// check for image state: Ready / Invalid / Client timeout
		log.Println("Checking image states....")
		imgc, errc = db.AllImage(localDB)
		checkerRes := make(chan db.Image)
		checker := cloud.ImageStateChecker{
			Images:  imgc,
			Done:    done,
			Result:  checkerRes,
			Timeout: time.Minute,
		}
		checker.Run(thread)

		var okCnt, errCnt int
		for img := range checkerRes {
			err = localDB.Save(&img)
			if verbose {
				if img.Error != "" || img.State == "Invalid" {
					errCnt++
					log.Printf("Image upload error: %q\n", img.Error)
				} else {
					okCnt++
					log.Printf("Image %q is %q\n", img.Filename, img.State)
				}
			}
			if err != nil {
				panic(err)
			}
		}

		// check whether the read from local db failed
		if err = <-errc; err != nil {
			panic(err)
		}

		log.Printf("%d out of %d images are uploaded and ready.", okCnt, totalImg)
		if errCnt > 0 {
			log.Printf("%d images failed. Please try again later.", errCnt)
		}
		log.Printf("To inspect more, type: 'alti-cli myproj inspect -p %v'\n", id)

		// generate report of uploading
		if report != "" {
			log.Println("Generating csv upload report...")
			out, err := os.Create(report)
			if err != nil {
				panic(err)
			}
			defer out.Close()
			writer := csv.NewWriter(out)

			writer.Write([]string{"Filename", "State", "Error"})
			imgc, errc = db.AllImage(localDB)
			for img := range imgc {
				err = writer.Write([]string{img.Filename, img.State, img.Error})
				if err != nil {
					panic(err)
				}
			}
			writer.Flush()

			// check whether the read from local db failed
			if err = <-errc; err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	importCmd.AddCommand(importImageCmd)
	importImageCmd.Flags().StringVarP(&id, "id", "p", id, "Project id")
	importImageCmd.Flags().StringVarP(&dir, "dir", "d", dir, "Directory path")
	importImageCmd.Flags().StringVarP(&skip, "skip", "s", skip, "Regular expression to skip paths")
	importImageCmd.Flags().StringVarP(&report, "report", "r", dir, "Path of csv upload report output")
	importImageCmd.Flags().StringVarP(&method, "method", "m", dir, "Desired method of upload: 'direct', 's3' or 'oss'")
	importImageCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display individual image info")
	importImageCmd.Flags().IntVarP(&thread, "thread", "n", thread, "Number of threads to process, default is number of cores x 4")
	importImageCmd.MarkFlagRequired("id")
	importImageCmd.MarkFlagRequired("dir")
}
