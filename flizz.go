package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func main() {
	file, err := os.Open("test1.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decodedImg, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	// Create a new mutable RGBA image
	img := image.NewRGBA(decodedImg.Bounds())
	// Add this line to move your photo data into the mutable 'img'
	draw.Draw(img, img.Bounds(), decodedImg, decodedImg.Bounds().Min, draw.Src)

	BoxBlur(9, img)
}

func Kernal(img image.Image) color.RGBA {
	bounds := img.Bounds()
	kernalSize := img.Bounds().Dx()
	kernalTotalPixels := uint32(kernalSize * kernalSize)
	//fmt.Println(kernalSize, kernalTotalPixels)
	var sumR uint32
	var sumG uint32
	var sumB uint32
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			R, G, B, _ := img.At(x, y).RGBA()
			sumR += R
			sumG += G
			sumB += B
		}
	}

	sumR = sumR / kernalTotalPixels
	sumG = sumG / kernalTotalPixels
	sumB = sumB / kernalTotalPixels
	//fmt.Println(uint8(sumR>>8), uint8(sumG>>8), uint8(sumB>>8))
	return color.RGBA{uint8(sumR >> 8), uint8(sumG >> 8), uint8(sumB >> 8), 255}
}

func BoxBlur(kernalSize int, img *image.RGBA) {
	kernalSize = (kernalSize / 2) - 1
	bounds := img.Bounds() // the full bounds of the whole image to loop over
	//loop over every pixel in the image
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			rect := image.Rect(x-kernalSize, y-kernalSize, x+kernalSize+1, y+kernalSize+1)

			sub := img.SubImage(rect)

			dst.Set(x, y, Kernal(sub))
		}
	}

	// Save the output
	outFile, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	png.Encode(outFile, dst)
}
