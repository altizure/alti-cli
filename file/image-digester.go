package file

import (
	"path/filepath"
	"runtime"
	"sync"

	"github.com/jackytck/alti-cli/errors"
)

// ImageDigest is the product of reading a regular file in local file system.
type ImageDigest struct {
	IsImage  bool
	Path     string
	Filename string
	Filesize int64 // in bytes
	Width    int
	Height   int
	GP       float64
	SHA1     string
	Error    error
}

// ImageDigester reads path names from paths...
type ImageDigester struct {
	Done   <-chan struct{}
	Paths  <-chan string
	Result chan<- ImageDigest
}

// Digest reads path names from Paths and sends digests of the corresponding
// files on Result until either Paths or Done is closed.
func (id *ImageDigester) Digest() {
	for path := range id.Paths {
		select {
		case id.Result <- work(path):
		case <-id.Done:
			return
		}
	}
}

// Run starts n number of goroutines to digest image files.
// If n is not positive, it will be set to number of CPU cores.
// Return n.
func (id *ImageDigester) Run(n int) int {
	if n <= 0 {
		n = runtime.NumCPU()
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
func work(p string) ImageDigest {
	ret := ImageDigest{
		Path: p,
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

	// c. filesize
	bytes, err := Filesize(p)
	if err != nil {
		ret.Error = errors.ErrFilesize
		return ret
	}
	ret.Filesize = bytes

	// d. image dimension
	w, h, err := GetImageSize(p)
	if err != nil {
		ret.Error = errors.ErrFileImageDim
		return ret
	}
	ret.Width = w
	ret.Height = h

	// e. gp usage
	ret.GP = DimToGigaPixel(w, h)

	// f. checksum
	sha1, err := Sha1sum(p)
	if err != nil {
		ret.Error = errors.ErrFileChecksum
		return ret
	}
	ret.SHA1 = sha1

	return ret
}
