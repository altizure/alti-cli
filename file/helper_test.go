package file

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

// TestMain creates and removes nat.jpg and nat.png.
func TestMain(m *testing.M) {
	exitCode := run(m)
	os.Exit(exitCode)
}

func run(m *testing.M) int {
	const w, h = 1000, 1000
	var im draw.Image
	im = image.NewRGBA(image.Rectangle{Max: image.Point{X: w, Y: h}})
	im = fibGradient(im)
	f, err := os.Create("nat.jpg")
	if err != nil {
		panic(err)
	}
	err = jpeg.Encode(f, im, nil)
	if err != nil {
		panic(err)
	}
	f.Close()
	f, err = os.Create("nat.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, im)
	if err != nil {
		panic(err)
	}
	f.Close()

	// teardown
	defer func() {
		err = os.Remove("nat.jpg")
		if err != nil {
			panic(err)
		}
		err = os.Remove("nat.png")
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
		{"nonexisting", args{"nat"}, "", true},
		{"jpg", args{"nat.jpg"}, "image/jpeg", false},
		{"png", args{"nat.png"}, "image/png", false},
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
		{"nonexisting", args{"nat"}, "", true},
		{"jpg", args{"nat.jpg"}, "672ecff0c4cab64e77321c38a091b1a2fb3ede66", false},
		{"png", args{"nat.png"}, "3d92ce14e0d333df4df8c1b6adb922a6d5b3ecb3", false},
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
