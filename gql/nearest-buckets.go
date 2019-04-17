package gql

import (
	"context"
	"strings"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/machinebox/graphql"
)

// BucketSuggestion returns the nearest buckets suggested by api.
func BucketSuggestion(cloud string) (string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		{
		  getGeoIPInfo {
		    nearestBuckets {
		      cloud
		      bucket
		    }
		  }
		}
	`)
	req.Header.Set("key", active.Key)
	ctx := context.Background()

	var res nearestBucketsRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	buks := res.GetGeoIPInfo.NearestBuckets
	if len(buks) == 0 {
		return "", errors.ErrNoBucketSuggestion
	}

	for _, b := range buks {
		if strings.ToLower(b.Cloud) == strings.ToLower(cloud) {
			return b.Bucket, nil
		}
	}
	return "", errors.ErrNoBucketSuggestion
}

type nearestBucketsRes struct {
	GetGeoIPInfo struct {
		NearestBuckets []struct {
			Cloud  string
			Bucket string
		}
	}
}
