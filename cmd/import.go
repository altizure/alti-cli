package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Root command for all import related commands",
	Long:  `'alti-cli import image' to import the image(s) of a directory`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("See alti-cli help import")
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
