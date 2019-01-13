package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackytck/alti-cli/errors"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "alti-cli",
	Short: "An Altizure CLI",
	Long:  `A CLI tool for interacting with Altizure service.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.altizure/config)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		altiDir := filepath.Join(home, ".altizure")
		confPath := filepath.Join(altiDir, "config.yaml")
		if _, err := os.Stat(altiDir); os.IsNotExist(err) {
			err := os.Mkdir(altiDir, 0755)
			errors.Must(err)
		}
		if _, err := os.Stat(confPath); os.IsNotExist(err) {
			_, err = os.Create(filepath.Join(altiDir, "config.yaml"))
			errors.Must(err)
		}

		// Search config in home directory + ".altizure".
		viper.SetConfigType("yaml")
		viper.AddConfigPath(altiDir)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
