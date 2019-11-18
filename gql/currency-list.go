package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/text"
	"github.com/machinebox/graphql"
)

// QueryCurrency infers the exact currency name from query string currency.
func QueryCurrency(currency string) (string, []string, error) {
	list, err := CurrencyList()
	if err != nil {
		return "", list, err
	}
	ret := text.BestMatch(list, currency, "")
	if ret == "" {
		return ret, list, errors.ErrCurrencyInvalid
	}
	return ret, list, nil
}

// CurrencyList returns a list of available currency supported by the api server.
func CurrencyList() ([]string, error) {
	var ret []string

	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	req := graphql.NewRequest(`
		query ($type: String!) {
			__type(name: $type) {
				enumValues {
					name
				}
			}
		}
	`)

	req.Header.Set("key", active.Key)
	req.Var("type", "CURRENCY")

	ctx := context.Background()
	var res currencyRes
	if err := client.Run(ctx, req, &res); err != nil {
		return ret, err
	}

	for _, c := range res.Type.EnumValues {
		ret = append(ret, c.Name)
	}

	return ret, nil
}

type currencyRes struct {
	Type enumCurType `json:"__type"`
}

type enumCurType struct {
	EnumValues []struct {
		Name string
	}
}
