package cloud

import (
	"runtime"
	"sync"

	"github.com/jackytck/alti-cli/db"
	"github.com/jackytck/alti-cli/errors"
)

// ImageStateChecker check the image states of all images.
type ImageStateChecker struct {
	Images <-chan db.Image
	Done   <-chan struct{}
	Result chan<- db.Image
}

// Digest checks state of each image from Images and send back the
// result to Result until either Images or Done is closed.
func (isc *ImageStateChecker) Digest() {
	for img := range isc.Images {
		select {
		case isc.Result <- isc.checkState(img):
		case <-isc.Done:
			return
		}
	}
}

// Run starts n number of goroutines to digest each image.
// If n is not positive, it will be set to number of CPU cores.
// Return n.
func (isc *ImageStateChecker) Run(n int) int {
	if n <= 0 {
		n = runtime.NumCPU()
	}
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			isc.Digest()
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(isc.Result)
	}()

	return n
}

func (isc *ImageStateChecker) checkState(img db.Image) db.Image {
	var ret db.Image
	ret.Error = errors.ErrNotImplemented.Error()
	return ret
}
