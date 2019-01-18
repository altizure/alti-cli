package gql

import (
	"context"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// MySelf queries simple info of current user.
func MySelf() (string, *types.User, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

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
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res mySelfRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", nil, err
	}

	if res.My.Self.Email == "" {
		return "", nil, errors.ErrNotLogin
	}

	u := types.User{
		Email:    res.My.Self.Email,
		Name:     res.My.Self.Name,
		Username: res.My.Self.Username,
	}
	return active.Endpoint, &u, nil
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
