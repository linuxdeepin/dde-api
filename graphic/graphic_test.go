/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"testing"

	"pkg.deepin.io/lib/dbusutil"

	C "gopkg.in/check.v1"
	"pkg.deepin.io/lib/utils"
)

func Test(t *testing.T) { C.TestingT(t) }

var _ = C.Suite(&Graphic{})

const (
	originImg               = "testdata/origin_1920x1080.jpg"
	originImgWidth          = 1920
	originImgHeight         = 1080
	originImgDominantColorH = 198.6
	originImgDominantColorS = 0.40
	originImgDominantColorV = 0.43

	originImgPngSmall       = "testdata/origin_small_200x200.png"
	originImgPngSmallWidth  = 200
	originImgPngSmallHeight = 200
	originImgPngIcon1       = "testdata/origin_icon_1_48x48.png"
	originImgPngIcon2       = "testdata/origin_icon_2_48x48.png"
	originIconWidth         = 48
	originIconHeight        = 48
)

// data uri for originImgPngIcon2
const testDataUri = `data:image/png;base64,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`

func delta(x, y float64) float64 {
	if x >= y {
		return x - y
	}
	return y - x
}

func sumFileMd5(f string) (md5 string) {
	md5, _ = utils.SumFileMd5(f)
	return
}

func (g *Graphic) SetUpSuite(c *C.C) {
	var err error
	g.service, err = dbusutil.NewSessionService()
	if err != nil {
		c.Skip(fmt.Sprintf("failed to get service: %v", err))
	}
}

func (g *Graphic) TestBlurImage(c *C.C) {
	resultFile := "testdata/test_blurimage.png"
	err := g.BlurImage(originImg, resultFile, 50, 1, "png")
	if err != nil {
		c.Error(err)
	}
}

func (g *Graphic) TestClipImage(c *C.C) {
	resultFile := "testdata/test_clipimage_100x200.png"
	err := g.ClipImage(originImg, resultFile, 0, 0, 100, 200, "png")
	if err != nil {
		c.Error(err)
	}
}

func (g *Graphic) TestConvertImage(c *C.C) {
	resultFile := "testdata/test_convertimage.png"
	err := g.ConvertImage(originImg, resultFile, "png")
	if err != nil {
		c.Error(err)
	}
}

func (g *Graphic) TestConvertImageToDataUri(c *C.C) {
	dataUri, err := g.ConvertImageToDataUri(originImgPngIcon2)
	if err != nil {
		c.Error(err)
	}
	c.Check(dataUri, C.Equals, testDataUri)
}

func (g *Graphic) TestConvertDataUriToImage(c *C.C) {
	resultFile := "testdata/test_convertdatauri.png"
	err := g.ConvertDataUriToImage(testDataUri, resultFile, "png")
	if err != nil {
		c.Error(err)
	}
}

func (g *Graphic) TestCompositeImage(c *C.C) {
	resultFile := "testdata/test_compositeimage.png"
	err := g.CompositeImage(originImgPngSmall, originImgPngIcon1, resultFile, 0, 0, "png")
	if err != nil {
		c.Error(err)
	}
	err = g.CompositeImage(resultFile, originImgPngIcon2, resultFile, 24, 24, "png")
	if err != nil {
		c.Error(err)
	}
}

func (g *Graphic) TestCompositeImageUri(c *C.C) {
	resultFile := "testdata/test_compositeimageuri.png"
	srcImageUri, _ := g.ConvertImageToDataUri(originImgPngSmall)
	compImageUri1, _ := g.ConvertImageToDataUri(originImgPngIcon1)
	compImageUri2, _ := g.ConvertImageToDataUri(originImgPngIcon2)
	resultDataUri, _ := g.CompositeImageUri(srcImageUri, compImageUri1, 0, 0, "png")
	resultDataUri, _ = g.CompositeImageUri(resultDataUri, compImageUri2, 24, 24, "png")
	err := g.ConvertDataUriToImage(resultDataUri, resultFile, "png")
	if err != nil {
		c.Error(err)
	}
}

func (g *Graphic) TestGetDominantColorOfImage(c *C.C) {
	h, s, v, err := g.GetDominantColorOfImage(originImg)
	if err != nil {
		c.Error(err)
	}
	if delta(h, originImgDominantColorH) > 1 ||
		delta(s, originImgDominantColorS) > 0.1 ||
		delta(v, originImgDominantColorV) > 0.1 {
		c.Error("h, s, v = ", h, s, v)
	}
}

func (g *Graphic) TestFillImage(c *C.C) {
	resultFile := "testdata/test_flllimage_tile_200x200.png"
	err := g.FillImage(originImg, resultFile, 200, 200, "tile", "png")
	if err != nil {
		c.Error(err)
	}

	resultFile = "testdata/test_flllimage_tile_1600x1000.png"
	err = g.FillImage(originImg, resultFile, 1600, 1000, "tile", "png")
	if err != nil {
		c.Error(err)
	}

	resultFile = "testdata/test_flllimage_center_400x400.png"
	err = g.FillImage(originImg, resultFile, 400, 400, "center", "png")
	if err != nil {
		c.Error(err)
	}

	resultFile = "testdata/test_flllimage_center_1600x1000.png"
	err = g.FillImage(originImg, resultFile, 1600, 1000, "center", "png")
	if err != nil {
		c.Error(err)
	}
}

func (g *Graphic) TestFlipImageHorizontal(c *C.C) {
	resultFile := "testdata/test_flipimagehorizontal.png"
	err := g.FlipImageHorizontal(originImg, resultFile, "png")
	if err != nil {
		c.Error(err)
	}
}

func (g *Graphic) TestFlipImageVertical(c *C.C) {
	resultFile := "testdata/test_flipimagevertical.png"
	err := g.FlipImageVertical(originImg, resultFile, "png")
	if err != nil {
		c.Error(err)
	}
}

func (g *Graphic) TestHsv(c *C.C) {
	for r := 0; r < 255; r += 7 {
		for gg := 0; gg < 255; gg += 5 {
			for b := 0; b < 255; b += 3 {
				r0, g0, b0 := uint8(r), uint8(gg), uint8(b)
				h, s, v, _ := g.Rgb2Hsv(r0, g0, b0)
				r1, g1, b1, _ := g.Hsv2Rgb(h, s, v)
				if delta(float64(r0), float64(r1)) > 1 || delta(float64(g0), float64(g1)) > 1 || delta(float64(b0), float64(b1)) > 1 {
					c.Fatalf("r0, g0, b0 = %d, %d, %d   r1, g1, b1 = %d, %d, %d", r0, g0, b0, r1, g1, b1)
				}
			}
		}
	}
}

func (g *Graphic) TestGetImageSize(c *C.C) {
	w, h, err := g.GetImageSize(originImg)
	if err != nil {
		c.Error(err)
	}
	c.Check(int(w), C.Equals, originImgWidth)
	c.Check(int(h), C.Equals, originImgHeight)
}

func (g *Graphic) TestResizeImage(c *C.C) {
	resultFile := "testdata/test_scaleimage_500x600.png"
	err := g.ResizeImage(originImg, resultFile, 500, 600, "png")
	if err != nil {
		c.Error(err)
	}
	w, h, err := g.GetImageSize(resultFile)
	if err != nil {
		c.Error(err)
	}
	c.Check(int(w), C.Equals, 500)
	c.Check(int(h), C.Equals, 600)
}

func (g *Graphic) TestThumbnailImage(c *C.C) {
	resultFile := "testdata/test_thumbnail.png"
	maxWidth, maxHeight := 200, 200
	err := g.ThumbnailImage(originImg, resultFile, uint32(maxWidth), uint32(maxHeight), "png")
	if err != nil {
		c.Error(err)
	}
	w, h, _ := g.GetImageSize(resultFile)
	c.Check(int(w) <= maxWidth, C.Equals, true)
	c.Check(int(h) <= maxHeight, C.Equals, true)
}

func (g *Graphic) TestRotateImageLeft(c *C.C) {
	resultFile := "testdata/test_rotateimageleft.png"
	err := g.RotateImageLeft(originImg, resultFile, "png")
	if err != nil {
		c.Error(err)
	}
}

func (g *Graphic) TestRotateImageRight(c *C.C) {
	resultFile := "testdata/test_rotateimageright.png"
	err := g.RotateImageRight(originImg, resultFile, "png")
	if err != nil {
		c.Error(err)
	}
}
