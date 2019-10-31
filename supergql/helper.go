package supergql

import (
	"github.com/jackytck/alti-cli/config"
	"github.com/machinebox/graphql"
)

// SuperRequest returns the super gql request and client.
func SuperRequest(gql string) (*graphql.Request, *graphql.Client) {
	c := config.Load().GetActive()
	client := graphql.NewClient(c.Endpoint + "/super")

	req := graphql.NewRequest(gql)
	req.Header.Set("key", c.Key)
	req.Header.Set("altitoken", c.Token)

	return req, client
}
