package cloud

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	tmpDir       string // for storing newly created multipart files
}

// Run starts the registration and uploading process.
// Return the state of imported model.
func (mru *ModelRegUploader) Run() (string, error) {
	switch mru.Method {
	case service.DirectUploadMethod:
		return mru.directUpload()
	case service.S3UploadMethod:
		fallthrough
	case service.MinioUploadMethod:
		return mru.smUpload(mru.Method)
	}
	return "", errors.ErrUploadMethodInvalid
}

// Done cleanups this uploader if user wants to terminate early.
func (mru *ModelRegUploader) Done() {
	if mru.tmpDir != "" {
		err := os.RemoveAll(mru.tmpDir)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Removed %q\n", mru.tmpDir)
	}
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

	// c. signal completing upload
	state, err := gql.DoneModelUpload(mru.PID, false)
	if err != nil {
		return state, err
	}

	return mru.checkState()
}

// smUpload uploads to s3 or minio via a single zip or multipart way of uploading.
func (mru *ModelRegUploader) smUpload(method string) (string, error) {
	if mru.MultipartDir != "" {
		return mru.smUploadMulti7z(method)
	}
	return mru.smUploadSingle(method)
}

// smUploadMulti uploads each multipart of a obj zip to s3 or minio.
// Each part could be concatenated in raw binary form.
func (mru *ModelRegUploader) smUploadMulti(method string) (string, error) {
	tmpDir, err := ioutil.TempDir(".", "")
	if err != nil {
		return "", err
	}
	mru.tmpDir = tmpDir
	log.Printf("Created dir %q for storing parts\n", tmpDir)
	defer mru.Done()

	parts, err := file.SplitFile(mru.ModelPath, tmpDir, 0, mru.Verbose)
	if err != nil {
		return "", err
	}

	// a. register each part
	err = mru.uploadParts(method, tmpDir, parts, true)
	if err != nil {
		return "", err
	}

	// c. signal completing multipart upload
	return gql.DoneModelUpload(mru.PID, true)
}

// smUploadMulti7z uploads 7z multipart to s3 or minio.
func (mru *ModelRegUploader) smUploadMulti7z(method string) (string, error) {
	files, err := ioutil.ReadDir(mru.MultipartDir)
	if err != nil {
		return "", err
	}

	var parts []string
	for _, f := range files {
		parts = append(parts, f.Name())
	}
	err = mru.uploadParts(method, mru.MultipartDir, parts, false)
	if err != nil {
		return "", err
	}

	// c. signal completing multipart upload
	return gql.DoneModelUpload(mru.PID, false)
}

// uploadParts registers and uploads each part to s3.
// baseDir is the dir that contains all the parts.
// parts is the slice of filenames of each part.
func (mru *ModelRegUploader) uploadParts(method string, baseDir string, parts []string, removePart bool) error {
	for _, p := range parts {
		if mru.Verbose {
			log.Printf("Uploading %q\n", p)
		}
		localPath := filepath.Join(baseDir, p)
		var url string
		var err error
		switch method {
		case service.S3UploadMethod:
			_, url, err = gql.RegisterModelS3(mru.PID, mru.Bucket, p)
		case service.MinioUploadMethod:
			_, url, err = gql.RegisterModelMinio(mru.PID, mru.Bucket, p)
		}
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
				log.Printf("Retrying (x %d) upload to %s for %q\n", i+1, strings.Title(method), p)
			}
			time.Sleep(time.Second)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// smUploadSingle uploads a single obj zip to s3 or minio.
func (mru *ModelRegUploader) smUploadSingle(method string) (string, error) {
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
	if size > 5*1024 {
		log.Printf("Filesize (%.2f MB) is bigger than 5GB", size)
		return mru.smUploadMulti(method)
	}

	var url string
	switch method {
	case service.S3UploadMethod:
		_, url, err = gql.RegisterModelS3(mru.PID, mru.Bucket, mru.Filename)
	case service.MinioUploadMethod:
		_, url, err = gql.RegisterModelMinio(mru.PID, mru.Bucket, mru.Filename)
	}
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
			log.Printf("Retrying (x %d) upload to %s for %q\n", i+1, strings.Title(method), mru.Filename)
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return "", err
	}

	_, err = gql.DoneModelUpload(mru.PID, false)
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
