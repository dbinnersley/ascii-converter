package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	_ "golang.org/x/image/tiff"

	"github.com/nfnt/resize"
)

//ASCIIArt structure containing AsciiArt data
type ASCIIArt struct {
	Width  int
	Height int
	Values []string //Contains Individual characters as elements
}

//CreateLines Outputs an Array of ASCIIArt with each row as a new string
func (art *ASCIIArt) CreateLines() []string {
	output := make([]string, art.Height)
	for row := 0; row < art.Height; row++ {
		currentvalues := art.Values[row*art.Width : (row+1)*art.Width]
		stringrow := strings.Join(currentvalues, " ")
		output[row] = stringrow
	}
	return output
}

//ImageToASCII Converts the Image to Ascii
func ImageToASCII(im image.Image, width int, charSet string) *ASCIIArt {
	greyscale := ImageToGrey(im)
	resized := ResizeGray(greyscale, width)
	result := GrayToASCII(resized, charSet)
	return result
}

//ImageToGrey Converts an Image to Grayscale
func ImageToGrey(im image.Image) *image.Gray {
	bounds := im.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	outputImage := image.NewGray(image.Rect(0, 0, width, height))
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			outputImage.Set(i, j, im.At(i, j))
		}
	}
	return outputImage
}

//ResizeGray Resizes a gray image to the given width and
func ResizeGray(input *image.Gray, width int) *image.Gray {
	newWidth := uint(width)

	bounds := input.Bounds()
	xsize := bounds.Dx()
	ysize := bounds.Dy()
	proportion := float32(newWidth) / float32(xsize)
	newHeight := uint(float32(ysize) * proportion)

	resizedGray := resize.Resize(newWidth, newHeight, input, resize.Bilinear)
	return resizedGray.(*image.Gray)
}

//grayMinMax Gets the min and max values of the pixels in the gray image
func grayMinMax(im *image.Gray) [2]uint8 {
	maxval := uint8(0)
	minval := uint8(255)
	for _, value := range im.Pix {
		if value < minval {
			minval = value
		}
		if value > maxval {
			maxval = value
		}
	}
	return [2]uint8{minval, maxval}
}

//buildLUT Builds the lookup table mapping pixel intensity to ASCII Character
func buildLUT(charset string, minMax [2]uint8) [256]string {
	mappingcount := len(charset)
	bitmapping := strings.Split(charset, "")

	lut := [256]string{}
	valuerange := minMax[1] - minMax[0]
	if valuerange == 0 {
		valuerange = 1
	}

	for i := 0; i < 256; i++ {
		index := (((i - int(minMax[0])) * mappingcount) / int(valuerange))
		if index < 0 {
			index = 0
		}
		if index > mappingcount-1 {
			index = mappingcount - 1
		}
		lut[i] = bitmapping[index]
	}
	return lut
}

//GrayToASCII Converts a greyscale image to Ascii art
func GrayToASCII(im *image.Gray, charSet string) *ASCIIArt {
	bounds := im.Bounds()
	xmax := bounds.Dx()
	ymax := bounds.Dy()

	minMax := grayMinMax(im)
	lut := buildLUT(charSet, minMax)

	asciiValues := make([]string, ymax*xmax)
	for i := 0; i < len(im.Pix); i++ {
		asciiValues[i] = lut[im.Pix[i]]
	}

	result := &ASCIIArt{Width: xmax, Height: ymax, Values: asciiValues}
	return result
}
