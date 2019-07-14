package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/file"
	"github.com/jackytck/alti-cli/service"
	"github.com/spf13/cobra"
)

var outDir string
var chunkSize int64 = 100

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split a file into chunks",
	Long:  "Split a file into simple binary chunks. The parts are named as XXX.part.YYY.",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			if verbose {
				elapsed := time.Since(start)
				log.Println("Took", elapsed)
			}
		}()

		// ensure output dir
		if _, err := os.Stat(outDir); os.IsNotExist(err) {
			err := os.Mkdir(outDir, 0755)
			errors.Must(err)
		}

		// pre-checks general
		if err := service.Check(
			nil,
			service.CheckFile(model),
			service.CheckDir(outDir),
		); err != nil {
			log.Println(err)
			return
		}

		// warn if too many parts
		filesize, _ := file.Filesize(model)
		chunkSize = chunkSize * (1 << 20) // 2^20, MB to B
		numParts := filesize / int64(chunkSize)
		if numParts > 100 {
			var ans string
			fmt.Printf("Continue to split into %d parts or not? (Y/N): ", numParts)
			fmt.Scanln(&ans)
			ans = strings.ToUpper(ans)
			if ans != "Y" && ans != service.Yes {
				log.Println("Cancelled.")
				return
			}
		}

		// split
		parts, err := file.SplitFile(model, outDir, chunkSize, verbose)
		errors.Must(err)

		p := "part"
		if len(parts) > 1 {
			p += "s"
		}
		log.Printf("Written %d %s.\n", len(parts), p)
	},
}

func init() {
	fileCmd.AddCommand(splitCmd)
	splitCmd.Flags().StringVarP(&model, "file", "f", model, "File path of input file.")
	splitCmd.Flags().StringVarP(&outDir, "out", "o", outDir, "File path of output dir.")
	splitCmd.Flags().Int64VarP(&chunkSize, "size", "s", chunkSize, "Chunk size in mega bytes. Default to 100 MB.")
	splitCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display more info of operation")
	splitCmd.MarkFlagRequired("file")
	splitCmd.MarkFlagRequired("out")
}
