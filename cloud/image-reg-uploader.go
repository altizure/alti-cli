package cloud

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/jackytck/alti-cli/db"
	"github.com/jackytck/alti-cli/gql"
)

// ImageRegUploader coordinates image registration and uploading concurrently.
type ImageRegUploader struct {
	Method  string
	BaseURL string
	Images  <-chan db.Image
	Done    <-chan struct{}
	Result  chan<- string
}

// Digest registers and uploads each image from Images and send back the
// result to Result until either Images or Done is closed.
func (iru *ImageRegUploader) Digest() {
	for img := range iru.Images {
		select {
		case iru.Result <- fmt.Sprintf("TODO: %s", img.Filename):
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

func upload(method string, imgs []db.Image, baseURL string) error {
	log.Println("imgs", len(imgs), imgs[0].SID, imgs[0].Filename, imgs[len(imgs)-1].SID, imgs[len(imgs)-1].Filename)
	switch method {
	case "direct":
		log.Println("TODO: direct upload...")
		img := imgs[0]
		gqlImg, err := gql.RegisterImageURL(img.PID, baseURL+img.URL, img.Filename, img.Hash)
		if err != nil {
			return err
		}
		fmt.Println(gqlImg)
		for {
			gqlImg, err = gql.ProjectImage(img.PID, gqlImg.ID)
			if err != nil {
				return err
			}
			fmt.Println(gqlImg)
			if gqlImg.State != "Uploaded" {
				break
			}
			time.Sleep(time.Second * 1)
		}
	case "s3":
		log.Println("TODO: s3 upload...")
	case "oss":
		log.Println("TODO: oss upload...")
	}

	return nil
}
