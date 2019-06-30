package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackytck/alti-cli/gql"
)

// SuggestUploadMethod suggests the best upload method if it is not set.
// Prefer direct upload over s3 over oss.
// kind is "image" or "model".
// Return "direct", "s3", "oss", ""
func SuggestUploadMethod(method, kind string) string {
	if method != "" {
		return strings.ToLower(method)
	}

	silent := func(string, ...interface{}) {}

	// check direct upload
	err := CheckDirectUpload(false, silent)
	if err == nil {
		return "direct"
	}

	// check s3
	sups := gql.SupportedCloud("", "", kind)
	var hasS3, hasOSS bool
	for _, s := range sups {
		if s == "S3" {
			hasS3 = true
		}
		if s == "OSS" {
			hasOSS = true
		}
	}

	if hasS3 {
		return "s3"
	}

	if hasOSS {
		return "oss"
	}

	return ""
}

// SuggestBucket suggests the best bucket if it is not set.
// And check if the bucket is valid if is set.
// Prefer the geo closest and supported one.
// kind is "image", "model" or "meta".
func SuggestBucket(method, bucket, kind string) (string, error) {
	if method == DirectUploadMethod {
		return "", nil
	}
	if bucket == "" {
		b, err := gql.SuggestedBucket(kind, method)
		if err != nil {
			return "", err
		}
		return b, nil
	}

	b, buckets, err := gql.QueryBucket(kind, method, bucket)
	if err != nil {
		e := fmt.Sprintf("Valid buckets are: %q\n", buckets)
		return "", errors.New(e)
	}
	return b, nil
}
