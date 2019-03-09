package cmd

import (
	"log"
	"time"

	"github.com/jackytck/alti-cli/file"
	"github.com/spf13/cobra"
)

var dir string
var skip string
var verbose bool

// checkImageCmd represents the checkImage command
var checkImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Check images of given directory recursively",
	Long: `Compute checksum, find duplicates and compute total giga-pixel
of all images of a given directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		log.Printf("Checking %s...\n", dir)

		var totalGP float64
		var totalImg int
		var totalMB float64

		done := make(chan struct{})
		defer close(done)

		paths, errc := file.WalkFiles(done, dir, skip)
		result := make(chan file.ImageDigest)

		digester := file.ImageDigester{
			Done:   done,
			Paths:  paths,
			Result: result,
		}
		threads := digester.Run(-1)
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
				log.Printf("Path: %q, Filename: %q, Dimension: %dx%d, GP: %.2f, Size: %.2f MB, Checksum: %s\n",
					r.Path, r.Filename, r.Width, r.Height, r.GP, mb, r.SHA1)
			}

			totalGP += r.GP
			totalImg++
			totalMB += mb
		}

		// check whether the Walk failed
		if err := <-errc; err != nil {
			panic(err)
		}

		if totalImg > 0 {
			log.Printf("Found %d images, total %.2f GP", totalImg, totalGP)
		} else {
			log.Println("No image is found!")
		}

		if verbose {
			log.Println("Took", time.Since(start))
		}
	},
}

func init() {
	checkCmd.AddCommand(checkImageCmd)
	checkImageCmd.Flags().StringVarP(&dir, "dir", "d", dir, "Directory path")
	checkImageCmd.Flags().StringVarP(&skip, "skip", "s", skip, "Regular expression to skip paths")
	checkImageCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display individual image info")
	checkImageCmd.MarkFlagRequired("dir")
}
