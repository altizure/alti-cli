package cmd

import (
	"fmt"
	"log"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/spf13/cobra"
)

var cash float64
var currency = "USD"

// toCoinsCmd converts cash to coins
var toCoinsCmd = &cobra.Command{
	Use:   "tocoins",
	Short: "Convert cash to coins",
	Long:  "Convert cash to coins in different currencies",
	Run: func(cmd *cobra.Command, args []string) {
		cur, err := service.SuggestCurrency(currency)
		if err != nil {
			log.Println(err)
			return
		}
		coins, err := gql.MoneyToCoins(cash, cur)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("%s%.2f could buy %.2f coins\n", cur, cash, coins)
	},
}

func init() {
	bankCmd.AddCommand(toCoinsCmd)
	toCoinsCmd.Flags().Float64VarP(&cash, "cash", "c", cash, "Number of cash to convert")
	toCoinsCmd.Flags().StringVarP(&currency, "currency", "f", currency, "Type of currency to convert (default USD)")
	errors.Must(toCoinsCmd.MarkFlagRequired("cash"))
}
