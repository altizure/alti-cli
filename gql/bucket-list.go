package gql

import (
	"strings"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/text"
)

var bucketType = map[string]map[string]string{
	"image": {
		"s3":    "BucketS3",
		"oss":   "BucketOSS",
		"minio": "BucketMinio",
	},
	"model": {
		"s3":    "BucketS3Model",
		"minio": "BucketMinioModel",
	},
	"meta": {
		"s3":    "BucketS3",
		"minio": "BucketMinioMeta",
	},
}

// QueryBucket infers the exact bucket name from query string bucket.
// kind is "image", "model" or "meta".
// cloud is "s3", "oss" or "minio".
func QueryBucket(kind, cloud, bucket string) (string, []string, error) {
	list, err := BucketList(kind, cloud)
	if err != nil {
		return "", list, err
	}
	ret := text.BestMatch(list, bucket, "")
	if ret == "" {
		return ret, list, errors.ErrBucketInvalid
	}
	return ret, list, nil
}

// BucketList returns a list of available buckets supported by the api server.
// kind is "image", "model" or "meta".
// cloud is "s3", "oss" or "minio".
func BucketList(kind, cloud string) ([]string, error) {
	var ret []string

	kind = strings.ToLower(kind)
	cloud = strings.ToLower(cloud)
	b, ok := bucketType[kind]
	if !ok {
		return ret, errors.ErrBucketInvalid
	}
	t, ok := b[cloud]
	if !ok {
		return ret, errors.ErrBucketInvalid
	}

	return EnumValues(t)
}
