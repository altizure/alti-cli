package cloud

import (
	"net/http"
	"os"

	"github.com/jackytck/alti-cli/file"
)

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
