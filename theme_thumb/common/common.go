package common

import (
	"image"
	"image/draw"
	"image/png"
	"os"
)

func CompositeIcons(images []image.Image, width, height, iconSize, padding int) image.Image {
	iconNum := len(images)
	destImg := image.NewRGBA(image.Rect(0, 0, width, height))
	if iconNum == 0 {
		return destImg
	}

	y := (height - iconSize) / 2
	spaceW := width - iconSize*iconNum
	x := (spaceW - (iconNum-1)*padding) / 2

	for _, srcImg := range images {
		draw.Draw(destImg, image.Rect(x, y, x+iconSize, y+iconSize), srcImg, image.Pt(0, 0), draw.Src)
		x += iconSize + padding
	}

	return destImg
}

func SavePngFile(m image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	return png.Encode(f, m)
}
