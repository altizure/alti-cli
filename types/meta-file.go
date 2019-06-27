package types

import "time"

// MetaFile represents the gql 'ProjectMetaFile' type.
type MetaFile struct {
	ID       string
	State    string
	Name     string
	Filename string
	Filesize float64
	Date     time.Time
	Chucksum string
	Error    []string
}
