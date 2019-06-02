package cloud

import (
	"log"
	"time"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/file"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
)

// ModelRegUploader coordinates model registration and uploading.
type ModelRegUploader struct {
	Method    string
	PID       string
	ModelPath string
	Filename  string
	DirectURL string
	Timeout   int
	Verbose   bool
}

// Run starts the registration and uploading process.
// Return the state of imported model.
func (mru *ModelRegUploader) Run() (string, error) {
	switch mru.Method {
	case service.DirectUploadMethod:
		return mru.directUpload()
	case "s3":
		return mru.s3Upload()
	}
	return "", errors.ErrUploadMethodInvalid
}

// directUpload registers the model via direct upload method and query its state
// change until timeout. Return the state of project.
func (mru *ModelRegUploader) directUpload() (string, error) {
	// sha1sum
	checksum, err := mru.checksum()
	if err != nil {
		return "", err
	}

	// register model
	im, err := gql.RegisterModelURL(mru.PID, mru.DirectURL, mru.Filename, checksum)
	if err != nil {
		return "", err
	}
	log.Printf("Registered model with state: %q\n", im.State)

	return mru.checkState()
}

func (mru *ModelRegUploader) s3Upload() (string, error) {
	// @TODO
	return "", errors.ErrNotImplemented
}

// checkState checks if the model state is changed from Pending until timeout.
func (mru *ModelRegUploader) checkState() (string, error) {
	timeout, err := mru.directTimeout()
	if err != nil {
		return "", err
	}

	stateC := make(chan string)

	// check per second
	go func() {
		log.Println("Checking state...")
		for {
			p, err := gql.Project(mru.PID)
			errors.Must(err)
			s := p.ImportedState
			if s != service.Pending {
				stateC <- s
				return
			}
			time.Sleep(time.Second * 1)
		}
	}()

	var state string
	log.Printf("Client will be timeout in %d seconds\n", timeout)
	select {
	case <-time.After(time.Second * timeout):
		return service.Pending, errors.ErrClientTimeout
	case state = <-stateC:
		return state, nil
	}
}

// checksum computes the SHA1 sum.
func (mru *ModelRegUploader) checksum() (string, error) {
	if mru.Verbose {
		log.Printf("Computing checksum: %q...\n", mru.Filename)
	}
	checksum, err := file.Sha1sum(mru.ModelPath)
	if err != nil {
		return "", err
	}
	if mru.Verbose {
		log.Printf("SHA1: %s\n", checksum)
	}
	return checksum, nil
}

// filesize gives the size of model in MegaByte.
func (mru *ModelRegUploader) filesize() (float64, error) {
	size, err := file.Filesize(mru.ModelPath)
	if err != nil {
		return 0, err
	}
	mb := file.BytesToMB(size)
	if mru.Verbose {
		log.Printf("Size: %.2f MB\n", mb)
	}
	return mb, nil
}

// directTimeout estimates the timeout in seconds for checking the uploaded state.
func (mru *ModelRegUploader) directTimeout() (time.Duration, error) {
	if mru.Timeout > 0 {
		return time.Duration(mru.Timeout), nil
	}
	// assume speed is at least 1 MegaByte
	mb, err := mru.filesize()
	if err != nil {
		return 0, err
	}
	timeout := int(mb)
	if timeout < 60 {
		timeout = 60
	}
	return time.Duration(timeout), nil
}
