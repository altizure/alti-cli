package db

import (
	"path"

	"github.com/asdine/storm"
	"github.com/jackytck/alti-cli/config"
	"github.com/jackytck/alti-cli/rand"
)

// OpenDB opens a storm db from path.
func OpenDB(path string) (*storm.DB, error) {
	if path == "" {
		p, err := OpenPath()
		if err != nil {
			return nil, err
		}
		path = p
	}
	db, err := storm.Open(path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

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

	dbFile := path.Join(confDir, randStr+".db")
	return dbFile, nil
}
