package main

import (
	"image"
	"image/color"
	"os"
	"testing"
)

var minMaxFixture = [2]uint8{0, 255}

var charsetFixture = "*+-0"

var catImage image.Image
var gnomeImage image.Image
var marblesImage image.Image
var homerImage image.Image

func rGBImageFixture() *image.RGBA {
	rGBImage := image.NewRGBA(image.Rect(0, 0, 2, 2))
	rGBImage.Set(0, 0, color.RGBA{R: 10, G: 10, B: 10, A: 0})
	rGBImage.Set(0, 1, color.RGBA{R: 100, G: 100, B: 100, A: 0})
	rGBImage.Set(1, 0, color.RGBA{R: 50, G: 50, B: 50, A: 0})
	rGBImage.Set(1, 1, color.RGBA{R: 20, G: 20, B: 20, A: 0})
	return rGBImage
}

func bigRGBImageFixture() *image.RGBA {
	rGBImage := image.NewRGBA(image.Rect(0, 0, 2000, 2000))
	rGBImage.Set(0, 0, color.RGBA{R: 10, G: 10, B: 10, A: 0})
	rGBImage.Set(0, 1, color.RGBA{R: 100, G: 100, B: 100, A: 0})
	rGBImage.Set(1, 0, color.RGBA{R: 50, G: 50, B: 50, A: 0})
	rGBImage.Set(1, 1, color.RGBA{R: 20, G: 20, B: 20, A: 0})
	return rGBImage
}

func grayImageFixture() *image.Gray {
	grayImage := image.NewGray(image.Rect(0, 0, 4, 2))
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			grayImage.Set(i, j, color.Gray{Y: 10})
		}
	}
	grayImage.Set(1, 1, color.Gray{Y: 4})
	grayImage.Set(2, 1, color.Gray{Y: 155})
	return grayImage
}

func TestMain(m *testing.M) {

	openFile, _ := os.Open("test/cat.jpg")
	defer openFile.Close()
	catImage, _, _ = image.Decode(openFile)

	openFile, _ = os.Open("test/gnome.png")
	defer openFile.Close()
	gnomeImage, _, _ = image.Decode(openFile)

	openFile, _ = os.Open("test/marbles.tif")
	defer openFile.Close()
	marblesImage, _, _ = image.Decode(openFile)

	openFile, _ = os.Open("test/homer.gif")
	defer openFile.Close()
	homerImage, _, _ = image.Decode(openFile)

	os.Exit(m.Run())
}

func TestCreateLines(t *testing.T) {
	art := &ASCIIArt{
		Values: []string{"A", "B", "B", "A"},
		Width:  2,
		Height: 2,
	}
	output := art.CreateLines()
	truth := []string{"A B", "B A"}

	for i := 0; i < 2; i++ {
		if output[i] != truth[i] {
			t.Fatal("Output is not correct")
		}
	}
}

func TestImageToGrey(t *testing.T) {
	rGBImage := rGBImageFixture()
	grayImage := ImageToGrey(rGBImage)

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if grayImage.At(j, i).(color.Gray).Y != color.GrayModel.Convert(rGBImage.At(j, i)).(color.Gray).Y {
				t.Error("Grayscale conversion failed")
			}
		}
	}
}

func TestResizeGray(t *testing.T) {
	grayImage := grayImageFixture()
	newGrayImage := ResizeGray(grayImage, 2)

	newBounds := newGrayImage.Bounds()
	xsize := newBounds.Dx()
	ysize := newBounds.Dy()

	if xsize != 2 {
		t.Error("x dimensions were not scaled correctly")
	}

	if ysize != 1 {
		t.Error("y dimensions were not scaled correctly")
	}
}

func TestGrayMinMax(t *testing.T) {
	grayImage := grayImageFixture()

	minMax := grayMinMax(grayImage)

	if minMax[0] != 4 {
		t.Error("Minimum value was not found correctly")
	}

	if minMax[1] != 155 {
		t.Error("Maximum value was not found correctly")
	}
}

func TestBuildLUT(t *testing.T) {
	lut := buildLUT(charsetFixture, minMaxFixture)
	truth := [256]string{}
	for i := 0; i < 64; i++ {
		truth[i] = "*"
	}
	for i := 64; i < 128; i++ {
		truth[i] = "+"
	}
	for i := 128; i < 192; i++ {
		truth[i] = "-"
	}
	for i := 192; i < 256; i++ {
		truth[i] = "0"
	}
	for i := 0; i < 256; i++ {
		if lut[i] != truth[i] {
			t.Fatal("Lut values are not equal")
		}
	}

}

func TestGrayToASCII(t *testing.T) {

	grayImage := grayImageFixture()

	art := GrayToASCII(grayImage, charsetFixture)

	truth := []string{"*", "*", "*", "*", "*", "*", "0", "*"}

	for i := range truth {
		if truth[i] != art.Values[i] {
			t.Fatal("conversion from image to ascii failed")
		}
	}
}

func BenchmarkTifToASCII(b *testing.B) {
	ImageToASCII(marblesImage, 100, "@%#*+=-:. ")
}
func BenchmarkGifToASCII(b *testing.B) {
	ImageToASCII(homerImage, 100, "@%#*+=-:. ")
}
func BenchmarkPngToASCII(b *testing.B) {
	ImageToASCII(gnomeImage, 100, "@%#*+=-:. ")
}
func BenchmarkJpegToASCII(b *testing.B) {
	ImageToASCII(catImage, 100, "@%#*+=-:. ")
}
