package cmd

import (
	"log"
	"path/filepath"
	"time"

	"github.com/jackytck/alti-cli/file"
	"github.com/spf13/cobra"
)

var dir string
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
		pc := file.WalkDir(dir)

		var totalGP float64
		var totalImg int
		var totalMB float64

		for p := range pc {
			isImg, err := file.IsImageFile(p)
			if err != nil || !isImg {
				continue
			}
			bytes, err := file.Filesize(p)
			if err != nil {
				log.Printf("Error getting filesize of %s\n", p)
				continue
			}
			mb := file.BytesToMB(bytes)
			w, h, err := file.GetImageSize(p)
			if err != nil {
				log.Printf("Error getting dimension of %s\n", p)
				continue
			}
			gp := file.DimToGigaPixel(w, h)
			sha1, err := file.Sha1sum(p)
			if err != nil {
				log.Printf("Error getting sha1 checksum of %s\n", p)
				continue
			}
			filename := filepath.Base(p)

			if verbose {
				log.Printf("Filename: %q, Dimension: %dx%d, GP: %.2f, Size: %.2f MB, Checksum: %s\n", filename, w, h, gp, mb, sha1)
			}

			totalGP += gp
			totalImg++
			totalMB += mb
		}

		if totalImg > 0 {
			log.Printf("Found %d images, total %.2f GP", totalImg, totalGP)
		} else {
			log.Println("No image is found!")
		}

		log.Println("Took", time.Since(start))
	},
}

func init() {
	checkCmd.AddCommand(checkImageCmd)
	checkImageCmd.Flags().StringVarP(&dir, "dir", "d", dir, "Directory path")
	checkImageCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display individual image info")
	checkImageCmd.MarkFlagRequired("dir")
}
