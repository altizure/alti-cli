package gql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/errors"
	"github.com/machinebox/graphql"
)

// SuggestedBucket returns the nearest bucket from api server.
// kind is "image", "model" or "meta".
// cloud is "s3", "oss" or "minio".
func SuggestedBucket(kind, cloud string) (string, error) {
	config := config.Load()
	active := config.GetActive()
	client := graphql.NewClient(active.Endpoint + "/graphql")

	// make a request
	req := graphql.NewRequest(fmt.Sprintf(`
		{
			getGeoIPInfo {
				%s {
					cloud
					bucket
				}
			}
		}
	`, kindToQuery(kind)))
	req.Header.Set("key", active.Key)
	ctx := context.Background()

	var res nearBucketRes
	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}
	var buks []cloudBucket
	switch kind {
	case "image":
		buks = res.GetGeoIPInfo.NearestBuckets
	case "meta":
		buks = res.GetGeoIPInfo.NearestMetaBuckets
	case "model":
		buks = res.GetGeoIPInfo.NearestModelBuckets
	}
	if len(buks) == 0 {
		return "", errors.ErrNoBucketSuggestion
	}

	for _, b := range buks {
		if strings.EqualFold(b.Cloud, cloud) {
			return b.Bucket, nil
		}
	}
	return "", errors.ErrNoBucketSuggestion
}

// kindToQuery gives the sub-query of getGeoIPInfo of a given kind.
// kind: "image", "meta" or "model"
func kindToQuery(kind string) string {
	switch kind {
	case "image":
		return "nearestBuckets"
	case "meta":
		return "nearestMetaBuckets"
	case "model":
		return "nearestModelBuckets"
	}
	return ""
}

type nearBucketRes struct {
	GetGeoIPInfo struct {
		NearestBuckets      []cloudBucket
		NearestMetaBuckets  []cloudBucket
		NearestModelBuckets []cloudBucket
	}
}

type cloudBucket struct {
	Cloud  string
	Bucket string
}
