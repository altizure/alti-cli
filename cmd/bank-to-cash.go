package cmd

import (
	"fmt"
	"log"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/spf13/cobra"
)

var coins float64

// toCashCmd converts coins to cash
var toCashCmd = &cobra.Command{
	Use:   "tocash",
	Short: "Convert coins to cash",
	Long:  "Convert coins to cash in different currencies",
	Run: func(cmd *cobra.Command, args []string) {
		cash, err := gql.CoinsToMoney(coins, currency)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("%.2f coins could be bought by %s%.2f", coins, currency, cash)
	},
}

func init() {
	bankCmd.AddCommand(toCashCmd)
	toCashCmd.Flags().Float64VarP(&coins, "coins", "c", coins, "Number of coins to convert")
	toCashCmd.Flags().StringVarP(&currency, "currency", "f", currency, "Type of currency to convert (default USD)")
	errors.Must(toCashCmd.MarkFlagRequired("coins"))
}
