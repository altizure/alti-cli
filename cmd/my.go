package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// myCmd represents the my command
var myCmd = &cobra.Command{
	Use:   "my",
	Short: "User related commands",
	Long:  "User related commands, such as my detail, settings and membership etc.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("See alti-cli help my")
	},
}

func init() {
	rootCmd.AddCommand(myCmd)
}
