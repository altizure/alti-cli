package service

import (
	"log"
	"strings"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/text"
	"github.com/jackytck/alti-cli/web"
)

// CheckAPIServer checks if API server is in normal mode.
func CheckAPIServer(logger func(string, ...interface{})) error {
	if logger == nil {
		logger = log.Printf
	}
	mode := gql.ActiveSystemMode()
	if mode != NormalMode {
		logger("API server is in %q mode.\n", mode)
		logger("Nothing could be uploaded at the moment!\n")
		return errors.ErrOffline
	}
	return nil
}

// CheckUploadMethod checks if the supplied upload method is suppored.
// kind is 'image', 'model' or 'meta'
func CheckUploadMethod(kind string, method string, logger func(string, ...interface{})) error {
	if logger == nil {
		logger = log.Printf
	}
	if method == "" {
		logger("No upload method is provided.")
		return errors.ErrUploadMethodInvalid
	}
	supMethods := gql.SupportedCloud("", "", kind)
	method = strings.ToLower(method)

	// check direct upload
	if method == DirectUploadMethod {
		err := CheckDirectUpload(false, logger)
		if err != nil {
			logger("Supported upload methods are: %q!", supMethods)
			return err
		}
		return nil
	}

	// check s3 or oss
	if sm := text.BestMatch(supMethods, method, ""); sm == "" {
		logger("Upload method: %q is not supported!\n", method)
		m := len(supMethods)
		switch m {
		case 0:
			logger("No supported mehtod is found! You could only use 'direct' upload!")
		case 1:
			logger("Only %q upload is supported!", supMethods[0])
		default:
			logger("Supported upload methods are: %q!", supMethods)
		}
		return errors.ErrUploadMethodInvalid
	}
	return nil
}

// CheckDirectUpload checks if direct upload is supported.
func CheckDirectUpload(verbose bool, logger func(string, ...interface{})) error {
	if logger == nil {
		logger = log.Printf
	}

	logger("Checking direct upload...")
	pu, _, err := web.PreferedLocalURL(verbose)
	if err != nil {
		logger("Client is invisible. Direct upload is not supported!")
		return err
	}
	logger("Direct upload is supported over %q!\n", pu.Hostname())
	return nil
}
