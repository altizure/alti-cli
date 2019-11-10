package file

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math"
	"regexp"

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

// IsFileExist checks if file exists.
func IsFileExist(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}
	return true
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

// IsZipFile tells if the file is a zip file.
func IsZipFile(f string) (bool, error) {
	ext, err := GuessFileType(f)
	if err != nil {
		return false, err
	}
	if ext == "application/zip" {
		return true, nil
	}
	return false, nil
}

// GuessFileType guesses the type of file.
func GuessFileType(file string) (string, error) {
	buff := make([]byte, 512)
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Read(buff)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buff), nil
}

// Sha1sum computes the sha1sum of the file.
func Sha1sum(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
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

// WalkFiles starts a goroutine to walk the directory tree at root and send the
// path of each regular file on the string channel.  It sends the result of the
// walk on the error channel. If done is closed, walkFiles abandons its work.
// skip is a regular expression pattern used for skipping paths. Would not skip
// if it is an empty string.
func WalkFiles(done <-chan struct{}, root string, skip string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	r, err := regexp.Compile(skip)
	if err != nil {
		errc <- err
		return paths, errc
	}

	onWalk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if skip != "" && r.MatchString(path) {
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

// SplitFile splits the file into parts and put it in the outDir.
// Each part would have chunkSize number of bytes.
// If chunkSize is larger than filesize, do nothing.
// If chunkSize is non-positive, will reset to 100 MB.
// Return filenames of parts.
func SplitFile(file, outDir string, chunkSize int64, verbose bool) ([]string, error) {
	if chunkSize <= 0 {
		chunkSize = 100 * (1 << 20) // 100MB
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	baseName := filepath.Base(file)

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	var partNames []string
	size := stat.Size()
	if chunkSize > size || baseName == "" {
		return partNames, nil
	}

	totalParts := uint64(math.Ceil(float64(size) / float64(chunkSize)))
	digit := int(math.Ceil(math.Log10(float64(totalParts))))

	for i := uint64(0); i < totalParts; i++ {
		partSize := chunkSize
		if i == totalParts-1 {
			partSize = size % chunkSize
			// evenly divide exactly
			if partSize == 0 {
				partSize = chunkSize
			}
		}
		buf := make([]byte, partSize)
		_, err := f.Read(buf)
		if err != nil {
			return nil, err
		}

		partName := fmt.Sprintf("%s.part.%0*d", baseName, digit, i+1)
		partPath := fmt.Sprintf("%s/%s", outDir, partName)

		if verbose {
			log.Printf("Writing %q\n", partName)
		}

		err = ioutil.WriteFile(partPath, buf, 0644)
		if err != nil {
			return nil, err
		}
		partNames = append(partNames, partName)
	}

	return partNames, nil
}

// MergeFile merges file parts into one single binary.
// It returns the number of bytes written and an error, if any.
func MergeFile(parts []string, output string) (int, error) {
	var merged []byte
	for _, p := range parts {
		f, err := os.Open(p)
		if err != nil {
			return 0, err
		}
		stat, err := f.Stat()
		if err != nil {
			return 0, err
		}
		bytes := make([]byte, stat.Size())
		buffer := bufio.NewReader(f)
		n, err := buffer.Read(bytes)
		if err != nil {
			return 0, err
		}
		merged = append(merged, bytes[:n]...)
		f.Close()
	}

	// write back
	f, err := os.Create(output)
	if err != nil {
		return 0, err
	}
	n, err := f.Write(merged)
	if err != nil {
		return 0, err
	}

	return n, nil
}
