package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// superCmd represents the super command
var superCmd = &cobra.Command{
	Use:   "super",
	Short: "Superuser tools",
	Long:  "Root command of superuser tools.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("See alti-cli help super")
	},
}

func init() {
	rootCmd.AddCommand(superCmd)
}
