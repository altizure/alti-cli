package cmd

import (
	"fmt"

	"github.com/jackytck/alti-cli/config"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout current user",
	Long:  "Logout the current user by emptying the user token if found in config file.",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.Load()
		conf.Token = ""
		err := conf.Save()
		if err != nil {
			panic(err)
		}
		fmt.Println("You are logout!")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
