package db

// Image represents an image in the db.
type Image struct {
	SID       int    `storm:"id,increment"`
	PID       string `storm:"index"`
	IID       string `storm:"index"`
	Name      string
	Filename  string
	URL       string
	LocalPath string
	Hash      string `storm:"index"`
	State     string `storm:"index"`
	Width     int
	Height    int
	GP        float64
}
