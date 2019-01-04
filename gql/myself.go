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
		      email
		      name
					username
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
		Email:    res.My.Self.Email,
		Name:     res.My.Self.Name,
		Username: res.My.Self.Username,
	}
	return &u, nil
}

type mySelfRes struct {
	My struct {
		Self struct {
			Email    string
			Name     string
			Username string
		}
	}
}
