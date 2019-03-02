package db

import (
	"path"

	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/rand"
)

// OpenPath infers a random path under the config directory
// for storing a temporary db.
func OpenPath() (string, error) {
	confDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}

	randStr, err := rand.RememberToken()
	if err != nil {
		return "", err
	}

	dbFile := path.Join(confDir, randStr, ".db")
	return dbFile, nil
}
