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
	client := graphql.NewClient(active.Endpoint + "/graphql")

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
	req.Header.Set("key", active.Key)
	req.Header.Set("altitoken", active.Token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var res mySelfRes
	if err := client.Run(ctx, req, &res); err != nil {
		switch err.(type) {
		case *url.Error:
			return active.Endpoint, nil, errors.ErrOffline
		default:
			return active.Endpoint, nil, err
		}
	}

	if res.My.Self.Email == "" {
		return "", nil, errors.ErrNotLogin
	}

	return active.Endpoint, &res.My.Self, nil
}

type mySelfRes struct {
	My struct {
		Self types.User
	}
}
