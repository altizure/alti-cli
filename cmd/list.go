package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Root command for listing various resources or variables.",
	Long:  `'alti-cli list bucket' to list all available buckets`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("See alti-cli help list")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
