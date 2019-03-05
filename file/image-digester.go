package file

// ImageDigest is the product of reading a regular file in local file system.
type ImageDigest struct {
	IsImage  bool
	Filename string
	Filesize int64
	Width    int64
	Height   int64
	GP       float64
	SHA1     string
}

// ImageDigester reads path names from paths...
type ImageDigester struct {
	// @TODO
}
