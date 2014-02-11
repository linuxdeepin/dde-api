package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func (dimg *DImage) GetImageSize(imageFile string) (w, h int32, err error) {
	// open the image file
	fr, err := os.Open(imageFile)
	if err != nil {
		// logError(err.Error()) // TODO
		return
	}
	defer fr.Close()

	img, _, err := image.Decode(fr)
	if err != nil {
		// image format not support
		// logError(err.Error()) // TODO
		return
	}

	w = int32(img.Bounds().Max.X)
	h = int32(img.Bounds().Max.Y)
	return
}
