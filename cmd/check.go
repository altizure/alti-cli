package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Root command for all check related commands",
	Long:  `'alti-cli check image' to check the image of a directory`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("See alti-cli help check")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
