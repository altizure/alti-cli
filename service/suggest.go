package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackytck/alti-cli/gql"
)

// SuggestUploadMethod suggests the best upload method if it is not set.
// Prefer direct upload over s3 over oss.
// kind is "image" or "model" or "meta".
// Return suggested method: "direct", "s3", "oss", ""
// and if it is suggested or not. If this is false, further checking is needed.
func SuggestUploadMethod(method, kind string) (string, bool) {
	if method != "" {
		return strings.ToLower(method), false
	}

	// check direct upload
	err := CheckDirectUpload(false, nil)
	if err == nil {
		return "direct", true
	}

	// check s3
	sups := gql.SupportedCloud("", "", kind)
	var hasS3, hasMinio, hasOSS bool
	for _, s := range sups {
		if s == "S3" {
			hasS3 = true
		}
		if s == "MINIO" {
			hasMinio = true
		}
		if s == "OSS" {
			hasOSS = true
		}
	}

	if hasS3 {
		return "s3", true
	}

	if hasMinio {
		return "minio", true
	}

	if hasOSS {
		return "oss", true
	}

	return "", false
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
		e := fmt.Sprintf("Valid buckets are: %q. Your input: %q\n", buckets, bucket)
		return "", errors.New(e)
	}
	return b, nil
}

// SuggestCurrency suggests the best match currency.
func SuggestCurrency(currency string) (string, error) {
	if currency == "" {
		return "USD", nil
	}

	c, validList, err := gql.QueryCurrency(currency)
	if err != nil {
		e := fmt.Sprintf("Valid currency are: %q. Your input: %q", validList, currency)
		return "", errors.New(e)
	}
	return c, nil
}
