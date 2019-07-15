package file

import (
	"bufio"
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jackytck/alti-cli/rand"
)

const testImgDir = "test/data/imgs/"
const testImgWidth = 1000
const testImgHeight = 1000

// TestMain creates and removes nat.jpg and nat.png.
func TestMain(m *testing.M) {
	exitCode := run(m)
	os.Exit(exitCode)
}

func run(m *testing.M) int {
	// create test dirs
	os.MkdirAll("test/data/imgs", os.ModePerm)
	os.Mkdir("test/data/other", os.ModePerm)

	// create test imgs
	var im draw.Image = image.NewRGBA(image.Rectangle{Max: image.Point{X: testImgWidth, Y: testImgHeight}})
	im = fibGradient(im)
	f, err := os.Create(testImgDir + "nat.jpg")
	if err != nil {
		panic(err)
	}
	err = jpeg.Encode(f, im, nil)
	if err != nil {
		panic(err)
	}
	f.Close()
	f, err = os.Create(testImgDir + "nat.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, im)
	if err != nil {
		panic(err)
	}
	f.Close()

	// unrelated stuff
	f, err = os.Create("test/data/other/log.txt")
	if err != nil {
		panic(err)
	}
	f.WriteString("2.71828182845904523536028747135266249775724709369995957496696")
	f.Close()

	// teardown
	defer func() {
		err = os.RemoveAll("test")
		if err != nil {
			panic(err)
		}
	}()

	return m.Run()
}

func fibGradient(im draw.Image) draw.Image {
	min, max := im.Bounds().Min, im.Bounds().Max
	for x := min.X; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			im.Set(x, y, color.RGBA{
				R: uint8(fib(x) % 256),
				G: uint8(fib(y) % 256),
				B: uint8(fib(x+y) % 256),
				A: uint8((x + y) % 256),
			})
		}
	}
	return im
}

func fib(n int) int {
	if n == 0 {
		return 0
	}
	a, b := 0, 1
	for i := 1; i < n; i++ {
		a, b = b, a+b
	}
	return b
}

func TestIsImageFile(t *testing.T) {
	type args struct {
		img string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"non-existing", args{"nat"}, false, true},
		{"jpg", args{testImgDir + "nat.jpg"}, true, false},
		{"png", args{testImgDir + "nat.png"}, true, false},
		{"non-image", args{"test/data/other/log.txt"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsImageFile(tt.args.img)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsImageFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsImageFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGuessFileType(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"none-xisting", args{"nat"}, "", true},
		{"jpg", args{testImgDir + "nat.jpg"}, "image/jpeg", false},
		{"png", args{testImgDir + "nat.png"}, "image/png", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GuessFileType(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("GuessFileType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GuessFileType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSha1sum(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"no-nexisting", args{"nat"}, "", true},
		{"jpg", args{testImgDir + "nat.jpg"}, "672ecff0c4cab64e77321c38a091b1a2fb3ede66", false},
		{"png", args{testImgDir + "nat.png"}, "3d92ce14e0d333df4df8c1b6adb922a6d5b3ecb3", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sha1sum(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sha1sum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Sha1sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalkFiles(t *testing.T) {
	done := make(chan struct{})
	paths, errc := WalkFiles(done, testImgDir, "")
	var got []string
	for p := range paths {
		got = append(got, p)
	}
	if err := <-errc; err != nil {
		t.Errorf("WalkFiles() error: %v", err)
	}
	want := []string{testImgDir + "nat.jpg", testImgDir + "nat.png"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WalkFiles() = %v, want %v", got, want)
	}
}

func TestWalkFilesSkip(t *testing.T) {
	done := make(chan struct{})
	paths, errc := WalkFiles(done, testImgDir, "\\w*.png")
	var got []string
	for p := range paths {
		got = append(got, p)
	}
	if err := <-errc; err != nil {
		t.Errorf("WalkFiles() error: %v", err)
	}
	want := []string{testImgDir + "nat.jpg"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WalkFiles() = %v, want %v", got, want)
	}
}

func TestGetImageSize(t *testing.T) {
	type args struct {
		img string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   int
		wantErr bool
	}{
		{"no-nexisting", args{"nat"}, 0, 0, true},
		{"jpg", args{testImgDir + "nat.jpg"}, testImgWidth, testImgHeight, false},
		{"png", args{testImgDir + "nat.png"}, testImgWidth, testImgHeight, false},
		{"non-image", args{"test/data/other/log.txt"}, 0, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetImageSize(tt.args.img)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetImageSize() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetImageSize() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDimToGigaPixel(t *testing.T) {
	type args struct {
		w int
		h int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"small", args{0, 0}, 0.0020736},
		{"iPhone X", args{4032, 3024}, 0.012192768},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DimToGigaPixel(tt.args.w, tt.args.h); got != tt.want {
				t.Errorf("DimToGigaPixel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_max(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"simple", args{1, 3}, 3},
		{"simple", args{520, 3}, 520},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := max(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilesize(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"non-existing", args{"nat"}, 0, true},
		{"jpg", args{testImgDir + "nat.jpg"}, 415568, false},
		{"png", args{testImgDir + "nat.png"}, 1734577, false},
		{"non-image", args{"test/data/other/log.txt"}, 61, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Filesize(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Filesize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Filesize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesToMB(t *testing.T) {
	type args struct {
		bytes int64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"simple", args{1734577}, 1.654221534729004},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToMB(tt.args.bytes); got != tt.want {
				t.Errorf("BytesToMB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitMerge(t *testing.T) {
	size := 1024
	original, err := rand.Bytes(size)
	if err != nil {
		t.Error(err)
	}
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)

	var parts []string
	t.Run("Split", func(t *testing.T) {
		data := filepath.Join(tmpDir, "data")
		err = ioutil.WriteFile(data, original, 0644)
		if err != nil {
			t.Error(err)
		}

		parts, err = SplitFile(data, tmpDir, 10, false)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Merge", func(t *testing.T) {
		var paths []string
		for _, p := range parts {
			paths = append(paths, filepath.Join(tmpDir, p))
		}
		merged := filepath.Join(tmpDir, "merged")
		n, err := MergeFile(paths, merged)
		if err != nil {
			t.Error(err)
		}
		if n != size {
			t.Errorf("Number of bytes written is incorrect\nGot %v\nWant %d\n", n, size)
		}

		f, err := os.Open(merged)
		if err != nil {
			t.Error(err)
		}
		bs := make([]byte, n)
		buffer := bufio.NewReader(f)
		_, err = buffer.Read(bs)
		if err != nil {
			t.Error(err)
		}
		f.Close()

		equal := bytes.Equal(original, bs)
		if !equal {
			t.Errorf("Bytes are not equal\nGot %v\nWant %v", bs, original)
		}
	})
}
