package file

import (
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
)

// ImageDigest is the product of reading a regular file in local file system.
type ImageDigest struct {
	IsImage  bool
	Path     string
	URL      string // relative url from root
	Filename string
	Filetype string
	Filesize int64 // in bytes
	Width    int
	Height   int
	GP       float64
	SHA1     string
	Existed  bool // existed in altizure or not
	Error    error
}

// ImageDigester reads path names from paths...
type ImageDigester struct {
	Root   string
	PID    string
	Done   <-chan struct{}
	Paths  <-chan string
	Result chan<- ImageDigest
}

// Digest reads path names from Paths and sends digests of the corresponding
// files on Result until either Paths or Done is closed.
func (id *ImageDigester) Digest() {
	for path := range id.Paths {
		select {
		case id.Result <- work(id.PID, id.Root, path):
		case <-id.Done:
			return
		}
	}
}

// Run starts n number of goroutines to digest image files.
// If n is not positive, it will be set to number of CPU cores x 4.
// Return n.
func (id *ImageDigester) Run(n int) int {
	if n <= 0 {
		n = runtime.NumCPU() * 4
	}
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			id.Digest()
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(id.Result)
	}()

	return n
}

// work checks the specified image file
// and get its name, size, width, height, gp and sha1.
func work(pid, r, p string) ImageDigest {
	ret := ImageDigest{
		Path: p,
		URL:  strings.Replace(p[len(r):], " ", "%20", -1),
	}

	// a. is image?
	isImg, err := IsImageFile(p)
	if err != nil {
		ret.Error = err
		return ret
	}
	if !isImg {
		ret.IsImage = false
		ret.Error = errors.ErrFileNotImage
		return ret
	}
	ret.IsImage = true

	// b. filename
	ret.Filename = filepath.Base(p)

	// c. filetype
	t, err := GuessFileType(p)
	if err != nil {
		ret.Error = err
		return ret
	}
	ret.Filetype = t

	// d. filesize
	bytes, err := Filesize(p)
	if err != nil {
		ret.Error = errors.ErrFilesize
		return ret
	}
	ret.Filesize = bytes

	// e. image dimension
	w, h, err := GetImageSize(p)
	if err != nil {
		ret.Error = errors.ErrFileImageDim
		return ret
	}
	ret.Width = w
	ret.Height = h

	// f. gp usage
	ret.GP = DimToGigaPixel(w, h)

	// g. checksum
	sha1, err := Sha1sum(p)
	if err != nil {
		ret.Error = errors.ErrFileChecksum
		return ret
	}
	ret.SHA1 = sha1

	// h. check if already uploaded
	ret.Existed, err = gql.HasImage(pid, sha1)
	if err != nil {
		ret.Error = err
		return ret
	}

	return ret
}
