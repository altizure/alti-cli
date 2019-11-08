package cmd

import (
	"log"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// quickCmd represents the quickCmd command
var quickCmd = &cobra.Command{
	Use:   "quick",
	Short: "Create and upload from directory",
	Long:  "Create a reconstruction project and upload images to it in a single command.",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			if verbose {
				elapsed := time.Since(start)
				log.Println("Took", elapsed)
			}
		}()

		// 1. create project
		name = filepath.Base(dir)
		newReconCmd.Run(cmd, args)

		// 2. import image
		id = newPID
		assumeYes = true
		verbose = true
		importImageCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(quickCmd)
	quickCmd.Flags().StringVarP(&dir, "dir", "d", dir, "Directory path")
}
