package cmd

import (
	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/gql"
)

// IsSuperuser checks if the currently active profile is a super user.
func IsSuperuser() bool {
	config := config.Load()
	active := config.GetActive()
	return gql.IsSuper(active.Endpoint, active.Key, active.Token)
}
