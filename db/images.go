package db

// Image represents an image in the db.
type Image struct {
	PID       string `storm:"id"`
	IID       string `storm:"index"`
	Name      string
	Filename  string
	LocalPath string
	Hash      string `storm:"index"`
	State     string `storm:"index"`
	Width     int
	Height    int
	GP        float64
}
