package cloud

import (
	"log"
	"time"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/file"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
)

// MetaFileRegUploader coordinates single meta file registration and uploading.
type MetaFileRegUploader struct {
	Method    string
	PID       string
	MID       string
	MetaPath  string
	Filename  string
	DirectURL string
	Bucket    string
	Timeout   int
	Verbose   bool
	checksum  string
}

// Run starts the registration and uploading process.
// Return the state of imported meta file.
func (mru *MetaFileRegUploader) Run() (string, error) {
	// check existence
	exists, err := mru.isUploaded()
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.ErrMetaExisted
	}

	// upload
	switch mru.Method {
	case service.DirectUploadMethod:
		return mru.directUpload()
	case "s3":
		return mru.s3Upload()
	}
	return "", errors.ErrUploadMethodInvalid
}

// Done cleanups this uploader if user wants to terminate early.
func (mru *MetaFileRegUploader) Done() error {
	return nil
}

func (mru *MetaFileRegUploader) isUploaded() (bool, error) {
	hash, err := mru.computeChecksum()
	if err != nil {
		return false, err
	}
	mru.checksum = hash
	return gql.HasMetaFile(mru.PID, hash)
}

// directUpload registers the meta file via direct upload method and query its state
// change until timeout. Return the state of meta file.
func (mru *MetaFileRegUploader) directUpload() (string, error) {
	// register meta file
	mf, err := gql.RegisterMetaURL(mru.PID, mru.DirectURL, mru.Filename, mru.checksum)
	if err != nil {
		return "", err
	}
	log.Printf("Registered meta with state: %q\n", mf.State)

	return mru.checkState()
}

// s3Upload uploads to s3.
func (mru *MetaFileRegUploader) s3Upload() (string, error) {
	if mru.Verbose {
		log.Printf("Uploading %q\n", mru.Filename)
	}
	size, err := mru.filesize()
	if err != nil {
		return "", err
	}
	if mru.Verbose {
		log.Printf("Size: %.2f MB\n", size)
	}
	meta, url, err := gql.RegisterMetaFileS3(mru.PID, mru.Bucket, mru.Filename)
	if err != nil {
		return "", err
	}
	mru.MID = meta.ID

	// b. upload to s3 with retry
	trial := 5
	for i := 0; i < trial; i++ {
		err = PutS3(mru.MetaPath, url)
		if err == nil {
			break
		}
		if mru.Verbose {
			log.Printf("Retrying (x %d) upload to S3 for %q\n", i+1, mru.Filename)
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return "", err
	}

	return mru.checkState()
}

// checkState checks if the model state is changed from Pending until timeout.
func (mru *MetaFileRegUploader) checkState() (string, error) {
	timeout, err := mru.getTimeout()
	if err != nil {
		return "", err
	}

	stateC := make(chan string)
	var stateErr error

	// check per second
	go func() {
		log.Println("Checking state...")
		for {
			m, err := gql.ProjectMetaFile(mru.PID, mru.MID)
			if err != nil {
				stateErr = err
				stateC <- ""
				return
			}
			s := m.State
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
		if stateErr != nil {
			return state, stateErr
		}
		return state, nil
	}
}

// computeChecksum computes the SHA1 sum.
func (mru *MetaFileRegUploader) computeChecksum() (string, error) {
	if mru.Verbose {
		log.Printf("Computing checksum: %q...\n", mru.Filename)
	}
	checksum, err := file.Sha1sum(mru.MetaPath)
	if err != nil {
		return "", err
	}
	if mru.Verbose {
		log.Printf("SHA1: %s\n", checksum)
	}
	return checksum, nil
}

// filesize gives the size of model in MegaByte.
func (mru *MetaFileRegUploader) filesize() (float64, error) {
	size, err := file.Filesize(mru.MetaPath)
	if err != nil {
		return 0, err
	}
	mb := file.BytesToMB(size)
	return mb, nil
}

// getTimeout estimates the timeout in seconds for checking the uploaded state.
func (mru *MetaFileRegUploader) getTimeout() (time.Duration, error) {
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
