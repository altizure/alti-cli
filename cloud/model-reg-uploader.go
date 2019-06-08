package cloud

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/file"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
)

// ModelRegUploader coordinates model registration and uploading.
type ModelRegUploader struct {
	Method       string
	PID          string
	ModelPath    string
	Filename     string
	DirectURL    string
	Bucket       string
	MultipartDir string // dir storing the 7zip multiparts
	Timeout      int
	Verbose      bool
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

// s3Upload uploads to s3 via a single zip or multipart way of uploading.
func (mru *ModelRegUploader) s3Upload() (string, error) {
	if mru.MultipartDir != "" {
		// return mru.s3UploadMulti()
		return mru.s3UploadMulti7z()
	}
	return mru.s3UploadSingle()
}

// s3UploadMulti uploads each multipart of a obj zip to s3.
// Each part could be concatenated in raw binary form.
func (mru *ModelRegUploader) s3UploadMulti() (string, error) {
	parts, err := file.SplitFile(mru.ModelPath, "/tmp", 1024*200)
	if err != nil {
		return "", err
	}

	// a. register each part
	err = mru.uploadParts("/tmp", parts, false)
	if err != nil {
		return "", err
	}

	// c. signal completing multipart upload
	return gql.DoneModelUpload(mru.PID)
}

func (mru *ModelRegUploader) s3UploadMulti7z() (string, error) {
	files, err := ioutil.ReadDir(mru.MultipartDir)
	if err != nil {
		return "", err
	}

	var parts []string
	for _, f := range files {
		parts = append(parts, f.Name())
	}
	err = mru.uploadParts(mru.MultipartDir, parts, false)
	if err != nil {
		return "", err
	}

	// c. signal completing multipart upload
	return gql.DoneModelUpload(mru.PID)
}

// uploadParts registers and uploads each part to s3.
// baseDir is the dir that contains all the parts.
// parts is the slice of filenames of each part.
func (mru *ModelRegUploader) uploadParts(baseDir string, parts []string, removePart bool) error {
	for _, p := range parts {
		if mru.Verbose {
			log.Printf("Uploading %q\n", p)
		}
		localPath := filepath.Join(baseDir, p)
		_, url, err := gql.RegisterModelS3(mru.PID, mru.Bucket, p)
		if err != nil {
			return err
		}

		// b. upload to s3 with retry
		trial := 5
		for i := 0; i < trial; i++ {
			err = PutS3(localPath, url)
			if err == nil {
				if removePart {
					os.Remove(localPath)
				}
				break
			}
			if mru.Verbose {
				log.Printf("Retrying (x %d) upload to S3 for %q\n", i+1, p)
			}
			time.Sleep(time.Second)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// s3UploadSingle uploads a single obj zip to s3.
// @TODO: precheck if file is bigger than 5GB
func (mru *ModelRegUploader) s3UploadSingle() (string, error) {
	if mru.Verbose {
		log.Printf("Uploading %q\n", mru.Filename)
	}
	_, url, err := gql.RegisterModelS3(mru.PID, mru.Bucket, mru.Filename)
	if err != nil {
		return "", err
	}

	// b. upload to s3 with retry
	trial := 5
	for i := 0; i < trial; i++ {
		err = PutS3(mru.ModelPath, url)
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
func (mru *ModelRegUploader) checkState() (string, error) {
	timeout, err := mru.getTimeout()
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

// getTimeout estimates the timeout in seconds for checking the uploaded state.
func (mru *ModelRegUploader) getTimeout() (time.Duration, error) {
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
