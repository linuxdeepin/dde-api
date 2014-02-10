package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func (color *Color) GetDominantColorOfImage(imagePath string) (h, s, v float64) {
	var def_h, def_s, def_v float64 = 200, 0.5, 0.8 // default hsv

	// open the image file
	fr, err := os.Open(imagePath)
	if err != nil {
		log.Printf(err.Error()) // TODO
		return def_h, def_s, def_v
	}
	defer fr.Close()

	img, _, err := image.Decode(fr)
	if err != nil {
		log.Printf(err.Error()) // TODO
		return def_h, def_s, def_v
	}

	// loop all points in image
	var sum_r, sum_g, sum_b, count uint64
	mx := img.Bounds().Max.X
	my := img.Bounds().Max.Y
	count = uint64(mx * my)
	if count == 0 {
		return def_h, def_s, def_v
	}
	for x := 0; x <= mx; x++ {
		for y := 0; y <= my; y++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			sum_r += uint64(r)
			sum_g += uint64(g)
			sum_b += uint64(b)
		}
	}

	h, s, v = color.RGB2HSV(uint8(sum_r/count), uint8(sum_g/count), uint8(sum_b/count))
	log.Printf("h=%f, s=%f, v=%f", h, s, v) // TODO
	if s < 0.05 {
		return def_h, def_s, def_v
	}
	return h, def_s, def_v
}
