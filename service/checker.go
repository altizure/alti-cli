package service

import (
	"log"
	"strings"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/text"
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
func CheckUploadMethod(method string, logger func(string, ...interface{})) error {
	if logger == nil {
		logger = log.Printf
	}
	if method != "" && strings.ToLower(method) != DirectUploadMethod {
		supMethods := gql.SupportedCloud("", "")
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
	}
	return nil
}
