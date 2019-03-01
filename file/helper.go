package file

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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
