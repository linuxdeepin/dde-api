package main

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
)

// Converts from any recognized format to PNG.
func (dimg *DImage) ConvertToPNG(src, dest string) (err error) {
	sf, err := os.Open(src)
	if err != nil {
		return
	}
	defer sf.Close()
	df, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer df.Close()

	img, _, err := image.Decode(sf)
	if err != nil {
		return
	}
	return png.Encode(df, img)
}

// Clip any recognized format image and save to PNG.
func (dimg *DImage) ClipPNG(src, dest string, x0, y0, x1, y1 int32) (err error) {
	sf, err := os.Open(src)
	if err != nil {
		return
	}
	defer sf.Close()

	df, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer df.Close()

	imgSrc, _, err := image.Decode(sf)
	if err != nil {
		return
	}

	imgDest := image.NewRGBA(image.Rect(int(x0), int(y0), int(x1), int(y1)))
	draw.Draw(imgDest, imgDest.Bounds(), imgSrc, image.Point{0, 0}, draw.Src)
	return png.Encode(df, imgDest)
}
