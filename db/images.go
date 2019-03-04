package db

// Image represents an image in the db.
type Image struct {
	PID       string `storm:"index"`
	IID       string `storm:"id"`
	Name      string
	Filename  string
	LocalPath string
	Hash      string `storm:"index"`
	State     string `storm:"index"`
	Width     int
	Height    int
	GP        float64
}
