package cmd

import (
	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/gql"
)

// LoginHint is shown when user wants to perfom operation that requires user token.
const LoginHint = "You are not login in!\nLogin with 'alti-cli login' or\nSwith account with 'alti-cli account use XXX'"

// IsLogin determines if user has logined.
func IsLogin() bool {
	config := config.Load()
	active := config.GetActive()
	return active.Token != ""
}

// IsSuperuser checks if the currently active profile is a super user.
func IsSuperuser() bool {
	config := config.Load()
	active := config.GetActive()
	return gql.IsSuper(active.Endpoint, active.Key, active.Token)
}
