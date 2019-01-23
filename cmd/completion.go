package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run:

[Linux]:
  sudo yum install bash-completion -y
  echo "source <(alti-cli completion)" >> ~/.bashrc
[Mac]:
  brew install bash-completion@2
  alti-cli completion > $(brew --prefix)/etc/bash_completion.d/alti-cli

Then restart shell.
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
