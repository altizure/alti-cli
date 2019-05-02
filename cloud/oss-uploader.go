package cloud

import (
	"fmt"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/types"
)

// NewOSSUploader returns a new OSSUploader.
func NewOSSUploader(pid string, refresh func() (*types.STS, error)) (*OSSUploader, error) {
	ret := OSSUploader{
		PID:        pid,
		RefreshSTS: refresh,
	}
	// retrieve STS and setup OSS connection
	err := ret.Refresh()
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// OSSUploader takes care of uploading files to a specific loccation
// and refresh its credentials.
type OSSUploader struct {
	PID        string
	RefreshSTS func() (*types.STS, error)
	creds      *types.STS
	bucket     *oss.Bucket
	mutex      sync.Mutex
}

// Refresh refreshes its STS token.
func (ou *OSSUploader) Refresh() error {
	sts, err := ou.RefreshSTS()
	if err != nil {
		return err
	}
	ou.setCreds(sts)
	return ou.reconnect()
}

// PutFile puts a file under the project's write-only space in OSS.
func (ou *OSSUploader) PutFile(filepath, cloudPath string) error {
	if expired, err := ou.isExpired(); expired || err != nil {
		ou.Refresh()
	}
	if expired, err := ou.isExpired(); expired || err != nil {
		if err != nil {
			return err
		}
		return errors.ErrNOSTS
	}
	key := fmt.Sprintf("%s/%s", ou.PID, cloudPath)
	return ou.getBucket().PutObjectFromFile(key, filepath)
}

// getCreds gets the sts creds.
func (ou *OSSUploader) getCreds() *types.STS {
	ou.mutex.Lock()
	defer ou.mutex.Unlock()
	return ou.creds
}

// setCreds sets the sts creds.
func (ou *OSSUploader) setCreds(creds *types.STS) {
	ou.mutex.Lock()
	defer ou.mutex.Unlock()
	ou.creds = creds
}

// getBucket gets the bucket handler.
func (ou *OSSUploader) getBucket() *oss.Bucket {
	ou.mutex.Lock()
	defer ou.mutex.Unlock()
	return ou.bucket
}

// setBucket sets the bucket handler.
func (ou *OSSUploader) setBucket(bucket *oss.Bucket) {
	ou.mutex.Lock()
	defer ou.mutex.Unlock()
	ou.bucket = bucket
}

// reconnect setups new oss connection from creds.
func (ou *OSSUploader) reconnect() error {
	// setup new connection
	sts := ou.getCreds()
	c, err := oss.New(sts.Endpoint, sts.ID, sts.Secret, oss.SecurityToken(sts.Token))
	if err != nil {
		return err
	}

	// bucket handler
	b, err := c.Bucket(sts.Bucket)
	if err != nil {
		return err
	}
	ou.setBucket(b)
	return nil
}

// isExpired tells if the current STS has expired.
func (ou *OSSUploader) isExpired() (bool, error) {
	due, err := time.Parse("2006-01-02T15:04:05Z", ou.getCreds().Expire)
	if err != nil {
		return true, err
	}
	now := time.Now()
	return now.After(due), nil
}
