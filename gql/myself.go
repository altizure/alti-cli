package gql

import (
	"context"
	"errors"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// MySelf queries simple info of current user.
func MySelf() (*types.User, error) {
	config, err := config.Load()
	if err != nil {
		return nil, errors.New("not login")
	}

	client := graphql.NewClient(config.Endpoint)

	// make a request
	req := graphql.NewRequest(`
		query {
		  my {
		    self {
		      name
		      email
				}
	    }
		}
	`)
	req.Header.Set("key", config.Key)
	req.Header.Set("altitoken", config.Token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res mySelfRes
	if err := client.Run(ctx, req, &res); err != nil {
		return nil, err
	}

	u := types.User{
		Name:  res.My.Self.Name,
		Email: res.My.Self.Email,
	}
	return &u, nil
}

type mySelfRes struct {
	My struct {
		Self struct {
			Name  string
			Email string
		}
	}
}
