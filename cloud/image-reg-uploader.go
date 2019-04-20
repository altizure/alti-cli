package cloud

import (
	"log"
	"runtime"
	"sync"

	"github.com/jackytck/alti-cli/db"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
)

// ImageRegUploader coordinates image registration and uploading concurrently.
type ImageRegUploader struct {
	Method  string
	Bucket  string
	BaseURL string
	Images  <-chan db.Image
	Done    <-chan struct{}
	Result  chan<- db.Image
}

// Digest registers and uploads each image from Images and send back the
// result to Result until either Images or Done is closed.
func (iru *ImageRegUploader) Digest() {
	for img := range iru.Images {
		select {
		case iru.Result <- iru.regUpload(img):
		case <-iru.Done:
			return
		}
	}
}

// Run starts n number of goroutines to digest each image.
// If n is not positive, it will be set to number of CPU cores.
// Return n.
func (iru *ImageRegUploader) Run(n int) int {
	if n <= 0 {
		n = runtime.NumCPU()
	}
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			iru.Digest()
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(iru.Result)
	}()

	return n
}

func (iru *ImageRegUploader) regUpload(img db.Image) db.Image {
	var ret db.Image
	switch iru.Method {
	case "direct":
		return iru.directUpload(img)
	case "s3":
		return iru.s3Upload(img)
	case "oss":
		return iru.ossUpload(img)
	}
	return ret
}

func (iru *ImageRegUploader) directUpload(img db.Image) db.Image {
	gqlImg, err := gql.RegisterImageURL(img.PID, iru.BaseURL+img.URL, img.Filename, img.Hash)
	if err != nil {
		img.Error = err.Error()
		return img
	}
	img.IID = gqlImg.ID
	img.State = gqlImg.State
	return img
}

func (iru *ImageRegUploader) s3Upload(img db.Image) db.Image {
	gqlImg, url, err := gql.RegisterImageS3(img.PID, iru.Bucket, img.Filename, img.Filetype, img.Hash)
	if err != nil {
		img.Error = err.Error()
		return img
	}
	img.IID = gqlImg.ID
	img.State = gqlImg.State
	// @TODO: upload to S3
	log.Println(url)
	return img
}

func (iru *ImageRegUploader) ossUpload(img db.Image) db.Image {
	img.Error = errors.ErrNotImplemented.Error()
	return img
}
