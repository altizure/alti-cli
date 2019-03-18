package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackytck/alti-cli/file"
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

		log.Printf("Checking %s...\n", dir)

		var totalGP float64
		var totalImg int
		var totalMB float64

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
		}

		// check whether the Walk failed
		if err := <-errc; err != nil {
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
		fmt.Printf("Continue to upload %d image%s or not? (Y/N): ", totalImg, plural)
		fmt.Scanln(&ans)
		ans = strings.ToUpper(ans)
		if ans != "Y" && ans != "YES" {
			log.Println("Cancelled.")
			return
		}

		// setup direct upload
		method := "direct"
		_, err := web.PreferedLocalURL()
		if err != nil {
			log.Println("Client is invisible. Direct upload is not supported!")
			method = "s3"
		} else {
			log.Println("Direct upload is supported!")
		}

		switch method {
		case "direct":
			log.Println("TODO: direct upload...")
		case "s3":
			log.Println("TODO: s3 upload...")
		case "oss":
			log.Println("TODO: oss upload...")
		}
	},
}

func init() {
	importCmd.AddCommand(importImageCmd)
	importImageCmd.Flags().StringVarP(&dir, "dir", "d", dir, "Directory path")
	importImageCmd.Flags().StringVarP(&skip, "skip", "s", skip, "Regular expression to skip paths")
	importImageCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display individual image info")
	importImageCmd.Flags().IntVarP(&thread, "thread", "n", thread, "Number of threads to process, default is number of cores")
	importImageCmd.MarkFlagRequired("dir")
}
