package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/asdine/storm"
	"github.com/jackytck/alti-cli/db"
	"github.com/jackytck/alti-cli/file"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/web"
	"github.com/spf13/cobra"
)

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
		log.Printf("Importing to %q...\n", p.Name)

		log.Printf("Checking %s...\n", dir)

		// stats
		var totalGP float64
		var totalImg int
		var totalMB float64

		// setup image digester
		done := make(chan struct{})
		defer close(done)

		paths, errc := file.WalkFiles(done, dir, skip)
		result := make(chan file.ImageDigest)

		digester := file.ImageDigester{
			Root:   dir,
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
		defer func() {
			err := os.Remove(dbPath)
			if err != nil {
				panic(err)
			}
		}()

		for r := range result {
			if r.Error != nil {
				log.Printf("Invalid image: %q, Reason: %v", r.Path, r.Error)
				continue
			}

			mb := file.BytesToMB(r.Filesize)
			if verbose {
				log.Printf("Path: %q, URL: %q, Filename: %q, Dimension: %d x %d, GP: %.2f, Type: %s, Size: %.2f MB, Checksum: %s\n",
					r.Path, r.URL, r.Filename, r.Width, r.Height, r.GP, r.Filetype, mb, r.SHA1)
			}

			totalGP += r.GP
			totalImg++
			totalMB += mb

			img := db.Image{
				PID:       p.ID,
				Filename:  r.Filename,
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
			log.Println("No image is found!")
			return
		}

		// ask user to proceed or not
		var ans string
		log.Printf("Found %d images, total %.2f GP, %.2f MB", totalImg, totalGP, totalMB)
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
		method := "direct"
		var localSever *http.Server
		var port int
		pu, _, err := web.PreferedLocalURL()
		if err != nil {
			log.Println("Client is invisible. Direct upload is not supported!")
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
			log.Printf("Serving files at http://%s:%d\n", pu.Hostname(), port)

			defer func() {
				if err := localSever.Shutdown(context.TODO()); err != nil {
					panic(err)
				}
			}()
		}

		// read from local db
		var imgs []db.Image
		limit, skip := 10, 0
		for skip < totalImg {
			err = localDB.All(&imgs, storm.Limit(limit), storm.Skip(skip))
			if err != nil {
				panic(err)
			}
			upload(method, localDB, imgs)

			skip += limit
		}
	},
}

func upload(method string, db *storm.DB, imgs []db.Image) error {
	log.Println("imgs", len(imgs), imgs[0].SID, imgs[0].Filename, imgs[len(imgs)-1].SID, imgs[len(imgs)-1].Filename)
	switch method {
	case "direct":
		log.Println("TODO: direct upload...")
		fmt.Println(imgs[0])
		time.Sleep(time.Second * 100)
	case "s3":
		log.Println("TODO: s3 upload...")
	case "oss":
		log.Println("TODO: oss upload...")
	}

	return nil
}

func init() {
	importCmd.AddCommand(importImageCmd)
	importImageCmd.Flags().StringVarP(&id, "id", "p", id, "Project id")
	importImageCmd.Flags().StringVarP(&dir, "dir", "d", dir, "Directory path")
	importImageCmd.Flags().StringVarP(&skip, "skip", "s", skip, "Regular expression to skip paths")
	importImageCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display individual image info")
	importImageCmd.Flags().IntVarP(&thread, "thread", "n", thread, "Number of threads to process, default is number of cores")
	importImageCmd.MarkFlagRequired("id")
	importImageCmd.MarkFlagRequired("dir")
}
