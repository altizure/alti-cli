package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/jackytck/alti-cli/cloud"
	"github.com/jackytck/alti-cli/db"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/file"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/jackytck/alti-cli/types"
	"github.com/jackytck/alti-cli/web"
	"github.com/spf13/cobra"
)

const directUpload = "direct"

var method string
var bucket string
var report string
var assumeYes bool

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

		// pre-checks general
		method = service.SuggestUploadMethod(method, "image")
		if err := service.Check(
			nil,
			service.CheckAPIServer(),
			service.CheckUploadMethod("image", method, ip, port),
			service.CheckPID("image", id),
			service.CheckDir(dir),
		); err != nil {
			log.Println(err)
			return
		}

		// get pid
		p, _ := gql.SearchProjectID(id, true)

		// setup direct upload server
		var serDone func()
		var baseURL string
		if method == service.DirectUploadMethod {
			bu, done, err := web.StartLocalServer(dir, ip, port, false)
			errors.Must(err)
			defer done()
			serDone = done
			baseURL = bu
		}

		// set bucket
		if method == "s3" || method == "oss" {
			log.Printf("Using %s to upload\n", method)
			if bucket == "" {
				b, err2 := gql.SuggestedBucket("image", method)
				if err2 != nil {
					panic(err2)
				}
				bucket = b
			} else {
				b, buckets, err2 := gql.QueryBucket("image", method, bucket)
				if err2 != nil {
					log.Printf("Valid buckets are: %q\n", buckets)
					return
				}
				bucket = b
			}
			log.Printf("Bucket %q is chosen", bucket)
		}

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
			serDone()
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
		if assumeYes {
			fmt.Println("Yes")
		} else {
			fmt.Scanln(&ans)
			ans = strings.ToUpper(ans)
			if ans != "Y" && ans != "YES" {
				log.Println("Cancelled.")
				return
			}
		}

		// read from local db, register and upload
		imgc, errc := db.AllImage(localDB)
		ruRes := make(chan db.Image)
		ruDigester := cloud.ImageRegUploader{
			Method:  method,
			Bucket:  bucket,
			BaseURL: baseURL,
			Images:  imgc,
			Done:    done,
			Result:  ruRes,
			Verbose: verbose,
		}
		if method == "oss" {
			err2 := ruDigester.WithOSSUploader(p.ID)
			if err2 != nil {
				panic(err2)
			}
		}
		ruDigester.Run(thread)

		regFailCnt := 0
		for img := range ruRes {
			err = localDB.Save(&img)
			if verbose {
				if img.Error != "" {
					log.Printf("Registration failed: %q\n", img.Error)
					regFailCnt++
				} else {
					if method == "direct" {
						log.Printf("Registered %q\n", img.Filename)
					} else {
						log.Printf("Registered and uploaded %q\n", img.Filename)
					}
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
		if regFailCnt == totalImg {
			log.Println("You run out of luck! All images failed to register!")
			return
		}

		// check for image state: Ready / Invalid / Client timeout
		log.Println("Checking image states....")
		imgc, errc = db.AllImage(localDB)
		checkerRes := make(chan db.Image)
		checker := cloud.ImageStateChecker{
			Images:  imgc,
			Done:    done,
			Result:  checkerRes,
			Timeout: time.Minute * time.Duration(timeout),
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
	importImageCmd.Flags().StringVarP(&report, "report", "r", report, "Path of csv upload report output")
	importImageCmd.Flags().StringVarP(&method, "method", "m", method, "Desired method of upload: 'direct', 's3' or 'oss'")
	importImageCmd.Flags().IntVarP(&timeout, "timeout", "t", timeout, "Timeout of checking upload state in seconds")
	importImageCmd.Flags().StringVar(&ip, "ip", ip, "IP address of ad-hoc local server for direct upload.")
	importImageCmd.Flags().StringVar(&port, "port", port, "Port of ad-hoc local server for direct upload.")
	importImageCmd.Flags().StringVarP(&bucket, "bucket", "b", bucket, "Desired bucket to upload for method: 's3' or 'oss'")
	importImageCmd.Flags().BoolVarP(&assumeYes, "assumeyes", "y", assumeYes, "Assume yes; assume that the answer to any question which would be asked is yes")
	importImageCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display individual image info")
	importImageCmd.Flags().IntVarP(&thread, "thread", "n", thread, "Number of threads to process, default is number of cores x 4")
	importImageCmd.MarkFlagRequired("id")
	importImageCmd.MarkFlagRequired("dir")
}
