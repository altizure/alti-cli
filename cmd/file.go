package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "For file inspection and manipulation",
	Long:  "Root command for doing file inspection and manipulation.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("See alti-cli help file")
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
}
