package gql

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/TylerBrock/colorjson"
	"github.com/jackytck/alti-cli/config"
	"github.com/machinebox/graphql"
)

// ActiveClient constructs the gql client for the currently active profile.
// Return the gql client, endpint, key and token.
func ActiveClient(room string) (*graphql.Client, string, string, string) {
	if room == "" {
		room = "graphql"
	}

	config := config.Load()
	active := config.GetActive()
	endpoint := active.Endpoint
	key := active.Key
	token := active.Token

	url := fmt.Sprintf("%s/%s", endpoint, room)
	client := graphql.NewClient(url)

	return client, endpoint, key, token
}

// PrettyPrint prints a raw json string into an indented colored string.
func PrettyPrint(data []byte) (string, error) {
	var obj map[string]interface{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return "", err
	}
	f := colorjson.NewFormatter()
	f.Indent = 2
	bs, err := f.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

// EnumValues gets the list of enum values by type name.
func EnumValues(typeName string) ([]string, error) {
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
	req.Var("type", typeName)

	ctx := context.Background()
	var res enumRes
	if err := client.Run(ctx, req, &res); err != nil {
		return ret, err
	}

	for _, t := range res.Type.EnumValues {
		ret = append(ret, t.Name)
	}
	sort.Strings(ret)

	return ret, nil
}

type enumRes struct {
	Type struct {
		EnumValues []struct {
			Name string `json:"name"`
		} `json:"enumValues"`
	} `json:"__type"`
}
