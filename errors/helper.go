package errors

import (
	"fmt"
)

// Must ensures error is nil, otherwise panic.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// MustGQL handles gql errors.
func MustGQL(err error, endpoint string) string {
	if err == nil {
		return ""
	}
	switch err {
	case ErrNoConfig:
		return fmt.Sprintf("Config not found.\nLogin with 'alti-cli login'")
	case ErrNotLogin:
		return fmt.Sprintf("You are not login in!\nLogin with 'alti-cli login' or\nSwith account with 'alti-cli account use XXX'")
	case ErrOffline:
		if endpoint == "" {
			endpoint = "Endpoint"
		}
		return fmt.Sprintf("%s is offline\nCheck status with 'alti-cli account'", endpoint)
	default:
		panic(err)
	}
}
