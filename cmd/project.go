package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Root command for all project related commands",
	Long:  `'alti-cli project new' to create an empty project`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("See alti-cli help project")
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
