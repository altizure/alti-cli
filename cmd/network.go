package cmd

import (
	"context"
	"fmt"

	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/web"
	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Check if api server could reach this client",
	Long:  "Locally start a web server and check if the api server could reach this server.",
	Run: func(cmd *cobra.Command, args []string) {
		s := web.Server{Directory: "/tmp"}
		server, port, err := s.ServeStatic(false)
		if err != nil {
			panic(err)
		}

		ip, _ := web.GetOutboundIP()
		url := fmt.Sprintf("http://%v:%v", ip, port)

		res := gql.CheckDirectNetwork(url)
		fmt.Println("Support direct upload?", res)

		if err := server.Shutdown(context.TODO()); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)
}
