package cloud

import (
	"net/http"
	"os"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/file"
)

// PutS3 is a helper func to put to s3.
func PutS3(localPath, url string) error {
	res, err2 := PutFile(localPath, url)
	if err2 != nil {
		return err2
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.ErrS3Error
	}
	return nil
}

// PutFile puts the local file specified in filepath to the remote url
// via http PUT.
func PutFile(filepath string, url string) (*http.Response, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	stats, err := f.Stat()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", url, f)
	if err != nil {
		return nil, err
	}
	t, err := file.GuessFileType(filepath)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", t)
	req.ContentLength = stats.Size()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return res, err
	}

	return res, nil
}
