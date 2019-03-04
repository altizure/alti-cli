package file

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"image"
	// for image.DecodeConfig
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DimToGigaPixel computes the giga-pixel from width and height.
func DimToGigaPixel(w, h int) float64 {
	return float64(max(2073600, w*h)) / 1000000000
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// GetImageSize decodes the width and height of an image.
func GetImageSize(img string) (int, int, error) {
	valid, err := IsImageFile(img)
	if err != nil {
		return 0, 0, err
	}
	if !valid {
		return 0, 0, nil
	}
	f, err := os.Open(img)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	i, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}
	return i.Width, i.Height, nil
}

// IsImageFile tells if the file is an image.
func IsImageFile(img string) (bool, error) {
	ext, err := GuessFileType(img)
	if err != nil {
		return false, err
	}
	if strings.Contains(ext, "image/") {
		return true, nil
	}
	return false, nil
}

// GuessFileType guesses the type of file.
func GuessFileType(file string) (string, error) {
	buff := make([]byte, 512)
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return "", err
	}
	f.Read(buff)
	return http.DetectContentType(buff), nil
}

// Sha1sum computes the sha1sum of the file.
func Sha1sum(file string) (string, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return "", err
	}
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Filesize returns the filesize in bytes of a file.
func Filesize(file string) (int64, error) {
	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}

	bytes := stat.Size()
	return bytes, nil
}

// BytesToMB converts bytes to mega-bytes.
func BytesToMB(bytes int64) float64 {
	return float64(bytes) / 1024 / 1024
}

// WalkDir walks the given directory and
// return each path through the returned channel.
func WalkDir(root string) <-chan string {
	paths := make(chan string)

	onWalk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		paths <- path
		return nil
	}

	go func() {
		filepath.Walk(root, onWalk)
		close(paths)
	}()

	return paths
}

// WalkFiles starts a goroutine to walk the directory tree at root and send the
// path of each regular file on the string channel.  It sends the result of the
// walk on the error channel. If done is closed, walkFiles abandons its work.
func WalkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	onWalk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		select {
		case paths <- path:
		case <-done:
			return errors.New("walk canceled")
		}
		return nil
	}

	go func() {
		// Close the paths channel after Walk returns.
		defer close(paths)
		// No select needed for this send, since errc is buffered.
		errc <- filepath.Walk(root, onWalk)
	}()

	return paths, errc
}
