// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package blurimage

import (
	"image"
	"image/color"
	"os"
	"path"
	"runtime/debug"
	"sync"

	"github.com/disintegration/imaging"
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
	var mu sync.Mutex

	imaging.AdjustFunc(img, func(c color.NRGBA) color.NRGBA {
		brightness := 0.2126*float64(c.R) + 0.7152*float64(c.G) + 0.0722*float64(c.B)
		mu.Lock()
		totalBrightness += brightness
		pixCount++
		mu.Unlock()
		return c
	})

	averBrightness := totalBrightness / pixCount

	// assume that average brightness higher than 100 is too bright.
	return averBrightness > 100
}
