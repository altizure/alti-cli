package gql

import (
	"context"
	"net/url"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
	"github.com/machinebox/graphql"
)

// MySelf queries simple info of current user.
func MySelf() (string, *types.User, error) {
	config := config.Load()
	active := config.GetActive()
	return MySelfByKeyToken(active.Endpoint, active.Key, active.Token)
}

// MySelfByKeyToken queries simple info of a specific user.
func MySelfByKeyToken(endpoint, key, token string) (string, *types.User, error) {
	client := graphql.NewClient(endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		query {
			my {
				self {
					email
					name
					username
					country
					balance
					freeGPQuota
					membershipState
					membership {
						state
						planName
						period
						startDate
						endDate
						memberGPQuota
						coinPerGP
						coinPriceDiscount
						assetStorage
						visibility
						coupon {
							value
							repeat
							validMonth
						}
						modelPerProject
						collaboratorQuota
						forceWatermark
					}
					modelUsage
					developer {
						status
					}
					stats {
						star
						project
						planet
						follower
						following
					}
					joined
				}
			}
		}
	`)
	req.Header.Set("key", key)
	req.Header.Set("altitoken", token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res mySelfRes
	if err := client.Run(ctx, req, &res); err != nil {
		switch err.(type) {
		case *url.Error:
			return endpoint, nil, errors.ErrOffline
		default:
			return endpoint, nil, err
		}
	}

	if res.My.Self.Email == "" {
		return "", nil, errors.ErrNotLogin
	}

	return endpoint, &res.My.Self, nil
}

type mySelfRes struct {
	My struct {
		Self types.User
	}
}
