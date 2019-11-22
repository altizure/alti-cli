package cloud

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/jackytck/alti-cli/db"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/jackytck/alti-cli/types"
)

// ImageRegUploader coordinates image registration and uploading concurrently.
type ImageRegUploader struct {
	Method  string
	Bucket  string
	BaseURL string
	Images  <-chan db.Image
	Done    <-chan struct{}
	Result  chan<- db.Image
	Verbose bool
	ossUp   *OSSUploader
}

// WithOSSUploader setups an OSS uploader for current pid and bucket.
func (iru *ImageRegUploader) WithOSSUploader(pid string) error {
	up, err := NewOSSUploader(pid, gql.RefreshSTS(pid, iru.Bucket))
	if err != nil {
		return err
	}
	iru.ossUp = up
	return nil
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
	case service.DirectUploadMethod:
		return iru.directUpload(img)
	case service.S3UploadMethod:
		fallthrough
	case service.MinioUploadMethod:
		return iru.smUpload(iru.Method, img, 5)
	case service.OSSUploadMethod:
		return iru.ossUpload(img)
	}
	return ret
}

func (iru *ImageRegUploader) directUpload(img db.Image) db.Image {
	u := fmt.Sprintf("%s/%s", iru.BaseURL, img.URL)
	gqlImg, err := gql.RegisterImageURL(img.PID, u, img.Filename, img.Hash)
	if err != nil {
		img.Error = err.Error()
		return img
	}
	img.IID = gqlImg.ID
	img.State = gqlImg.State
	return img
}

// smUpload uploads to either s3 or minio.
// kind is "s3" or "minio"
func (iru *ImageRegUploader) smUpload(kind string, img db.Image, retry int) db.Image {
	// a. register image
	var gqlImg *types.Image
	var url string
	var err error

	switch kind {
	case service.S3UploadMethod:
		gqlImg, url, err = gql.RegisterImageS3(img.PID, iru.Bucket, img.Filename, img.Filetype, img.Hash)
	case service.MinioUploadMethod:
		gqlImg, url, err = gql.RegisterImageMinio(img.PID, iru.Bucket, img.Filename, img.Filetype, img.Hash)
	}
	if err != nil {
		img.Error = err.Error()
		return img
	}
	img.IID = gqlImg.ID
	img.State = gqlImg.State

	// b. signal the start of upload
	trial := retry
	for i := 0; i < trial; i++ {
		state, e := gql.StartImageUpload(img.IID)
		err = e
		if e == nil {
			img.State = state
			break
		}
		if iru.Verbose {
			log.Printf("Retrying (x %d) mutating state for %q\n", i+1, img.Filename)
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		img.Error = err.Error()
		return img
	}

	// c. upload to s3/minio with retry
	// helper func to put to s3/minio
	upload := func() error {
		if iru.Verbose {
			log.Printf("Uploading %q\n", img.Filename)
		}
		res, err2 := PutFile(img.LocalPath, url)
		if err2 != nil {
			return err2
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			switch kind {
			case service.S3UploadMethod:
				return errors.ErrS3Error
			case service.MinioUploadMethod:
				return errors.ErrMinioError
			}
		}
		return nil
	}

	for i := 0; i < trial; i++ {
		err = upload()
		if err == nil {
			break
		}
		if iru.Verbose {
			log.Printf("Retrying (x %d) upload to %s for %q\n", i+1, kind, img.Filename)
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		img.Error = err.Error()
	}

	return img
}

func (iru *ImageRegUploader) ossUpload(img db.Image) db.Image {
	if iru.ossUp == nil {
		img.Error = errors.ErrOSSUploaderNotFound.Error()
		return img
	}

	// a. register oss image
	gqlImg, err := gql.RegisterImageOSS(img.PID, iru.Bucket, img.Filename, img.Filetype, img.Hash)
	if err != nil {
		img.Error = err.Error()
		return img
	}
	img.IID = gqlImg.ID
	img.State = gqlImg.State

	// b. signal the start of upload
	trial := 5
	for i := 0; i < trial; i++ {
		state, e := gql.StartImageUpload(img.IID)
		err = e
		if e == nil {
			img.State = state
			break
		}
		if iru.Verbose {
			log.Printf("Retrying (x %d) mutating state for %q\n", i+1, img.Filename)
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		img.Error = err.Error()
		return img
	}

	// c. upload to oss with retry
	for i := 0; i < trial; i++ {
		if iru.Verbose {
			log.Printf("Uploading %q\n", img.Filename)
		}
		err = iru.ossUp.PutFile(img.LocalPath, gqlImg.Filename)
		if err == nil {
			break
		}
		if iru.Verbose {
			log.Printf("Retrying (x %d) upload to OSS for %q\n", i+1, img.Filename)
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		img.Error = err.Error()
	}

	// d. signal the end of upload
	state, err := gql.DoneImageUpload(img.IID)
	if err != nil {
		img.Error = err.Error()
		return img
	}
	img.State = state

	return img
}
