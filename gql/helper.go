package gql

import (
	"encoding/json"
	"fmt"

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
