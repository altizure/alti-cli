package types

import "time"

// DownloadsConnection represents the gql 'DownloadsConnection' type.
type DownloadsConnection struct {
	TotalCount int
	Edges      []struct {
		Node Downloadable
	}
}

// Downloadable represents the gql 'Downloadable' type.
type Downloadable struct {
	State string
	Name  string
	Size  int64
	Mtime time.Time
	Link  string
}
