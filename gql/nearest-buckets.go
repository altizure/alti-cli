package gql

import (
	"context"
	"strings"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/machinebox/graphql"
)

// SuggestedBucket returns the nearest bucket from api server.
// 'kind' is 'image' or 'model'.
// 'cloud' is 's3', 'oss' or 'minio'.
func SuggestedBucket(kind, cloud string) (string, error) {
	switch kind {
	case "image":
		return imageBucketSuggestion(cloud)
	case "model":
		return modelBucketSuggestion(cloud)
	}
	return "", errors.ErrNoBucketSuggestion
}

// imageBucketSuggestion returns the nearest image bucket suggested by api.
func imageBucketSuggestion(cloud string) (string, error) {
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

	var res nearImgBuckRes
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

// modelBucketSuggestion returns the nearest model bucket suggested by api.
func modelBucketSuggestion(cloud string) (string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(`
		{
		  getGeoIPInfo {
		    nearestModelBuckets {
		      cloud
		      bucket
		    }
		  }
		}
	`)
	req.Header.Set("key", active.Key)
	ctx := context.Background()

	var res nearModelBuckRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	buks := res.GetGeoIPInfo.NearestModelBuckets
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

type nearImgBuckRes struct {
	GetGeoIPInfo struct {
		NearestBuckets []struct {
			Cloud  string
			Bucket string
		}
	}
}

type nearModelBuckRes struct {
	GetGeoIPInfo struct {
		NearestModelBuckets []struct {
			Cloud  string
			Bucket string
		}
	}
}
