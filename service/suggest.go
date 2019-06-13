package service

import (
	"strings"

	"github.com/jackytck/alti-cli/gql"
)

// SuggestUploadMethod suggests the best upload method if it is not set.
// Prefer direct upload over s3 over oss.
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
