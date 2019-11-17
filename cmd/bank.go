package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// bankCmd represents the bank command
var bankCmd = &cobra.Command{
	Use:   "bank",
	Short: "Convert between coins and cash",
	Long:  "Convert between alti coins and cash in different currencies",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("See alti-cli help bank")
	},
}

func init() {
	rootCmd.AddCommand(bankCmd)
}
