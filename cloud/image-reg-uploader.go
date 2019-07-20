package cloud

import (
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

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
	case "direct":
		return iru.directUpload(img)
	case "s3":
		return iru.s3Upload(img)
	case "minio":
		return iru.minioUpload(img)
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
	// a. register s3 image
	gqlImg, url, err := gql.RegisterImageS3(img.PID, iru.Bucket, img.Filename, img.Filetype, img.Hash)
	if err != nil {
		img.Error = err.Error()
		return img
	}
	img.IID = gqlImg.ID
	img.State = gqlImg.State

	// b. signal the start of upload
	state, err := gql.StartImageUpload(img.IID)
	if err != nil {
		img.Error = err.Error()
		return img
	}
	img.State = state

	// c. upload to s3 with retry
	// helper func to put to s3
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
			return errors.ErrS3Error
		}
		return nil
	}

	trial := 5
	for i := 0; i < trial; i++ {
		err = upload()
		if err == nil {
			break
		}
		if iru.Verbose {
			log.Printf("Retrying (x %d) upload to S3 for %q\n", i+1, img.Filename)
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		img.Error = err.Error()
	}

	return img
}

func (iru *ImageRegUploader) minioUpload(img db.Image) db.Image {
	// a. register minio image
	gqlImg, url, err := gql.RegisterImageMinio(img.PID, iru.Bucket, img.Filename, img.Filetype, img.Hash)
	if err != nil {
		img.Error = err.Error()
		return img
	}
	img.IID = gqlImg.ID
	img.State = gqlImg.State

	// b. signal the start of upload
	state, err := gql.StartImageUpload(img.IID)
	if err != nil {
		img.Error = err.Error()
		return img
	}
	img.State = state

	// c. upload to minio with retry
	// helper func to put to minio
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
			return errors.ErrS3Error
		}
		return nil
	}

	trial := 5
	for i := 0; i < trial; i++ {
		err = upload()
		if err == nil {
			break
		}
		if iru.Verbose {
			log.Printf("Retrying (x %d) upload to Minio for %q\n", i+1, img.Filename)
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
	state, err := gql.StartImageUpload(img.IID)
	if err != nil {
		img.Error = err.Error()
		return img
	}
	img.State = state

	// c. upload to oss with retry
	trial := 5
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
	state, err = gql.DoneImageUpload(img.IID)
	if err != nil {
		img.Error = err.Error()
		return img
	}
	img.State = state

	return img
}
