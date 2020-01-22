package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/file"
	"github.com/jackytck/alti-cli/service"
	"github.com/jackytck/alti-cli/text"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// checkImageGroup represents the checkImageGroup command
var checkImageGroupCmd = &cobra.Command{
	Use:   "image-group",
	Short: "Check if images are present in fs but not in group.txt",
	Long:  "Chcek if extra images are present in the fs but not defined in group.txt.",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			if verbose {
				elapsed := time.Since(start)
				log.Println("Took", elapsed)
			}
		}()

		log.Printf("Checking %s...\n", dir)

		// pre-checks general
		groupPath := path.Join(dir, "group.txt")
		if err := service.Check(
			nil,
			service.CheckDir(dir),
			service.CheckFile(groupPath),
		); err != nil {
			log.Println(err)
			return
		}

		// a. read group.txt
		group, err := readGroupTxt(groupPath)
		errors.Must(err)

		// b. read images
		var undefined []file.ImageDigest

		done := make(chan struct{})
		defer close(done)

		paths, errc := file.WalkFiles(done, dir, skip)
		result := make(chan file.ImageDigest)

		digester := file.ImageDigester{
			Root:      dir,
			Done:      done,
			Paths:     paths,
			Result:    result,
			LightWork: true,
		}
		threads := digester.Run(thread)
		if verbose {
			log.Printf("Working in %d thread(s)...", threads)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Undefined Filename"})

		for r := range result {
			if !r.IsImage {
				continue
			}
			if _, ok := group[r.Filename]; ok {
				continue
			}

			if verbose {
				log.Printf("Path: %q, URL: %q, Filename: %q\n", r.Path, r.URL, r.Filename)
			}

			if printTable {
				r := []string{
					fmt.Sprintf("%q", r.Filename),
				}
				table.Append(r)
			}

			undefined = append(undefined, r)
		}

		// check whether the Walk failed
		if err := <-errc; err != nil {
			panic(err)
		}

		// c. display results
		totalImg := len(undefined)
		plural := ""
		if totalImg > 1 {
			plural = "s"
		}
		if totalImg > 0 {
			log.Printf("Found %d undefined image%s\n", totalImg, plural)
		} else {
			log.Println("No image outside group.txt is found!")
		}

		if printTable {
			table.SetFooter([]string{fmt.Sprintf("%d image%s", totalImg, plural)})
			table.Render()
		}

		// d. ask if continue to remove
		fmt.Printf("Continue to remove %d undefined image%s or not? (Y/N): ", totalImg, plural)
		var ans string
		if assumeYes {
			fmt.Println("Yes")
		} else {
			fmt.Scanln(&ans)
			ans = strings.ToUpper(ans)
			if ans != "Y" && ans != service.Yes {
				log.Println("Cancelled.")
				return
			}
		}

		// e. remove undefined images
		for _, img := range undefined {
			if verbose {
				fmt.Printf("Removing %q\n", img.Path)
			}
			err := os.Remove(img.Path)
			errors.Must(err)
		}

		log.Println("Done")
	},
}

// readGroupTxt reads all the lines and return a map of image names.
func readGroupTxt(p string) (map[string]bool, error) {
	var ret map[string]bool
	lines, err := file.ReadFile(p)
	if err != nil {
		return ret, err
	}
	var s []string
	for _, l := range lines {
		toks := strings.Split(l, " ")
		s = append(s, toks[0])
	}
	return text.SliceToMap(s, true), nil
}

func init() {
	checkCmd.AddCommand(checkImageGroupCmd)
	checkImageGroupCmd.Flags().StringVarP(&dir, "dir", "d", dir, "Directory path")
	checkImageGroupCmd.Flags().StringVarP(&skip, "skip", "s", skip, "Regular expression to skip paths")
	checkImageGroupCmd.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Display individual image info")
	checkImageGroupCmd.Flags().BoolVarP(&printTable, "table", "t", printTable, "Output all of the found images in table format")
	checkImageGroupCmd.Flags().IntVarP(&thread, "thread", "n", thread, "Number of threads to process, default is number of cores x 4")
	checkImageGroupCmd.Flags().BoolVarP(&assumeYes, "assumeyes", "y", assumeYes, "Assume yes; assume that the answer to any question which would be asked is yes")
	errors.Must(checkImageGroupCmd.MarkFlagRequired("dir"))
}
