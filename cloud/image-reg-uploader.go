package cloud

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/jackytck/alti-cli/db"
	"github.com/jackytck/alti-cli/gql"
)

// ImageRegUploadRes contains the result with error of the image registration
// and uploading operation.
type ImageRegUploadRes struct {
	db.Image
	Error error
}

// ImageRegUploader coordinates image registration and uploading concurrently.
type ImageRegUploader struct {
	Method  string
	BaseURL string
	Images  <-chan db.Image
	Done    <-chan struct{}
	Result  chan<- ImageRegUploadRes
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

func (iru *ImageRegUploader) regUpload(img db.Image) ImageRegUploadRes {
	var ret ImageRegUploadRes
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

func (iru *ImageRegUploader) directUpload(img db.Image) ImageRegUploadRes {
	ret := ImageRegUploadRes{Image: img}
	gqlImg, err := gql.RegisterImageURL(img.PID, iru.BaseURL+img.URL, img.Filename, img.Hash)
	if err != nil {
		ret.Error = err
		return ret
	}
	fmt.Println("gqlImg", gqlImg)
	ret.State = gqlImg.State
	return ret
}

func (iru *ImageRegUploader) s3Upload(img db.Image) ImageRegUploadRes {
	ret := ImageRegUploadRes{Image: img}
	log.Println("TODO: s3 upload...")
	return ret
}

func (iru *ImageRegUploader) ossUpload(img db.Image) ImageRegUploadRes {
	ret := ImageRegUploadRes{Image: img}
	log.Println("TODO: oss upload...")
	return ret
}
