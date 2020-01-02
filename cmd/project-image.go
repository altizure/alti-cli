package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/jackytck/alti-cli/types"
	"github.com/spf13/cobra"
)

var out string

// exportImageCmd represents the image command
var exportImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Export all images to csv",
	Long:  "Export all images of a project to a csv.",
	Run: func(cmd *cobra.Command, args []string) {
		// a. check
		if err := service.Check(
			nil,
			service.CheckAPIServer(),
			service.CheckPID("image", id),
		); err != nil {
			log.Println(err)
			return
		}
		first := 10
		imgs, page, total, err := allImages(first, "")
		errors.Must(err)
		if total == 0 {
			log.Println("No image is found! Bye.")
			return
		}

		// b. setup csv writer
		if out == "" {
			out = fmt.Sprintf("%s-images.csv", id)
		}
		o, err := os.Create(out)
		errors.Must(err)

		defer o.Close()
		writer := csv.NewWriter(o)
		err = writer.Write([]string{"Filename", "State", "URL"})
		errors.Must(err)

		// c. export
		var cnt int
		log.Printf("Exporting %d images...\n", total)
		printProgress(cnt, total)

		work := func() {
			c, err := writeCSV(writer, imgs)
			if err != nil {
				panic(err)
			}
			cnt += c
			printProgress(cnt, total)
		}

		// d. loop all images in batch, fetch `first` images at a time
		work()
		for page.HasNextPage {
			imgs, page, _, err = allImages(first, page.EndCursor)
			if err != nil {
				panic(err)
			}
			work()
		}

		log.Println("Done")
	},
}

func printProgress(work, total int) {
	log.Printf("========== %v/%v ==========\n", work, total)
}

func writeCSV(w *csv.Writer, imgs []types.ProjectImage) (int, error) {
	for _, img := range imgs {
		fields := []string{
			img.Name,
			img.State,
			img.URL,
		}
		if verbose {
			log.Println(fields)
		}
		err := w.Write(fields)
		if err != nil {
			return 0, err
		}
	}
	w.Flush()
	return len(imgs), nil
}

func allImages(first int, after string) ([]types.ProjectImage, *types.PageInfo, int, error) {
	imgs, page, total, err := gql.AllProjectImages(id, first, 0, "", after)
	if msg := errors.MustGQL(err, ""); msg != "" {
		fmt.Println(msg)
		return nil, nil, 0, err
	}
	return imgs, page, total, nil
}

func init() {
	projectCmd.AddCommand(exportImageCmd)
	exportImageCmd.Flags().StringVarP(&id, "id", "p", id, "Project id")
	exportImageCmd.Flags().StringVarP(&out, "out", "o", out, "Path of output csv")
	exportImageCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display individual image info")
}
