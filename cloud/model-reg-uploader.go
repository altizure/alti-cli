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
func (mru *ModelRegUploader) Run() error {
	switch mru.Method {
	case service.DirectUploadMethod:
		return mru.directUpload()
	case "s3":
		return mru.s3Upload()
	}
	return nil
}

func (mru *ModelRegUploader) directUpload() error {
	checksum, err := mru.checksum()
	if err != nil {
		return err
	}
	timeout, err := mru.directTimeout()
	if err != nil {
		return err
	}

	// register model
	im, err := gql.RegisterModelURL(mru.PID, mru.DirectURL, mru.Filename, checksum)
	if err != nil {
		return err
	}
	log.Printf("Registered model with state: %q\n", im.State)

	// @TODO: refactor state checking
	// check if ready
	stateC := make(chan string)

	go func() {
		log.Println("Checking state...")
		for {
			p, err := gql.Project(mru.PID)
			errors.Must(err)
			s := p.ImportedState
			if s != "Pending" {
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
		log.Printf("Client timeout.")
		return errors.ErrClientTimeout
	case state = <-stateC:
		log.Printf("Model is in state: %q\n", state)
	}

	return nil
}

func (mru *ModelRegUploader) s3Upload() error {
	// @TODO
	return errors.ErrNotImplemented
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
