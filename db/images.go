package db

// Image represents an image in the db.
type Image struct {
	SID       int `storm:"id,increment"`
	PID       string
	IID       string
	Name      string
	Filename  string `storm:"index"`
	Filetype  string
	URL       string
	LocalPath string
	Hash      string
	State     string
	Width     int
	Height    int
	GP        float64
	Error     string
}
