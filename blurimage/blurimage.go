/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package blurimage

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"os"
	"path"
	"runtime/debug"
)

func BlurImage(file string, sigma float64, dest string) error {
	img, err := imaging.Open(file)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Dir(dest), 0755)
	if err != nil {
		return err
	}

	defer debug.FreeOSMemory()

	blurredNRGB := imaging.Blur(img, sigma)

	var finalNRGB image.Image = blurredNRGB
	// need to darken the image if it's too bright
	if isTooBright(blurredNRGB) {
		finalNRGB = imaging.AdjustBrightness(blurredNRGB, -20)
	}

	return imaging.Save(finalNRGB, dest)
}

func isTooBright(img image.Image) bool {
	var pixCount float64 = 0
	var totalBrightness float64 = 0

	imaging.AdjustFunc(img, func(c color.NRGBA) color.NRGBA {
		brightness := 0.2126*float64(c.R) + 0.7152*float64(c.G) + 0.0722*float64(c.B)
		totalBrightness += brightness
		pixCount++

		return c
	})

	averBrightness := totalBrightness / pixCount

	// assume that average brightness higher than 100 is too bright.
	return averBrightness > 100
}
