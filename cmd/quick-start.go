package cmd

import (
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackytck/alti-cli/file"
	"github.com/jackytck/alti-cli/service"
	"github.com/spf13/cobra"
)

var inputPath string

// quickCmd represents the quickCmd command
var quickCmd = &cobra.Command{
	Use:   "quick",
	Short: "Create and upload from directory or model zip",
	Long:  "Create a reconstruction project from a directory of images or an imported project from a model zip.",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			if verbose {
				elapsed := time.Since(start)
				log.Println("Took", elapsed)
			}
		}()

		// pre-check
		if err := service.Check(
			nil,
			service.CheckAPIServer(),
			service.CheckFile(inputPath),
			service.CheckDirOrZip(inputPath),
		); err != nil {
			log.Println(err)
			return
		}

		// 1. determine project type
		recon := true
		isZip, _ := file.IsZipFile(inputPath)
		if isZip {
			recon = false
		}

		// 2. create project
		if name == "" {
			// base file name as default project name
			name = filepath.Base(inputPath)

			if isZip {
				name = strings.TrimSuffix(name, path.Ext(name))
			}
		}
		if recon {
			// reconstruction project
			newReconCmd.Run(cmd, args)

			// 3a. import image
			id = newPID
			dir = inputPath
			assumeYes = true
			importImageCmd.Run(cmd, args)

			// 3b. import meta file
			metafiles, _ := service.GetMetafilePaths(inputPath)
			if len(metafiles) > 0 {
				bucket = ""
				for i, f := range metafiles {
					log.Printf("Detected metafile(%d/%d): %q\n", i+1, len(metafiles), f)
					meta = f
					importMetaCmd.Run(cmd, args)
				}
			}

			// 4. start reconstruction task
			startReconCmd.Run(cmd, args)
		} else {
			// imported model (obj zip) project
			newModelCmd.Run(cmd, args)

			// 3b. import model
			id = newPID
			model = inputPath
			importModelCmd.Run(cmd, args)
		}
	},
}

func init() {
	rootCmd.AddCommand(quickCmd)
	quickCmd.Flags().StringVarP(&inputPath, "input", "i", inputPath, "Directory path or model zip file")
	quickCmd.Flags().StringVarP(&name, "name", "n", name, "Project name")
	quickCmd.Flags().StringVarP(&projType, "projectType", "p", projType, "free, pro")
	quickCmd.Flags().StringVarP(&modelType, "modelType", "m", modelType, "CAD, PHOTOGRAMMETRY, PTCLOUD")
	quickCmd.Flags().StringVarP(&skip, "skip", "s", skip, "Regular expression to skip paths")
	quickCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display more info of operation")
}
