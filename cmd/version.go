package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version info.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.0.0")
	},
}
