package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/machinebox/graphql"
)

// CoinsToMoney converts coins into real currency.
func CoinsToMoney(coins float64, currency string) (float64, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		query ($coins: Float, $currency: CURRENCY) {
			bank {
				coinsToMoney(coins: $coins, currency: $currency)
			}
		}
	`)
	req.Var("coins", coins)
	req.Var("currency", currency)
	req.Header.Set("key", active.Key)

	ctx := context.Background()
	var res bankRes
	if err := client.Run(ctx, req, &res); err != nil {
		return 0, err
	}
	return res.bank.coinsToMoney, nil
}

// MoneyToCoins converts money into coins.
func MoneyToCoins(money float64, currency string) (float64, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		query ($money: Float, $currency: CURRENCY) {
			bank {
				moneyToCoins(money: $money, currency: $currency)
			}
		}
	`)
	req.Var("money", money)
	req.Var("currency", currency)
	req.Header.Set("key", active.Key)

	ctx := context.Background()
	var res bankRes
	if err := client.Run(ctx, req, &res); err != nil {
		return 0, err
	}
	return res.bank.moneyToCoins, nil
}

type bankRes struct {
	bank struct {
		coinsToMoney float64
		moneyToCoins float64
	}
}
