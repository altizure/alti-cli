package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// projectNewCmd represents the project new sub-command
var projNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Root command for all project creation related commands",
	Long:  `'alti-cli project new recon' to create an empty reconstruction project`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("See alti-cli help project new")
	},
}

func init() {
	projectCmd.AddCommand(projNewCmd)
}
