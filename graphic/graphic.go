/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"pkg.linuxdeepin.com/lib/dbus"
	libgdkpixbuf "pkg.linuxdeepin.com/lib/gdkpixbuf"
	libgraphic "pkg.linuxdeepin.com/lib/graphic"
)

const (
	graphicDest = "com.deepin.api.Graphic"
	graphicPath = "/com/deepin/api/Graphic"
	graphicIfs  = "com.deepin.api.Graphic"
)

// Graphic is a dbus interface wrapper for pkg.linuxdeepin.com/lib/graphic.
type Graphic struct{}

// GetDBusInfo implement interface of dbus.DBusObject
func (graphic *Graphic) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		graphicDest,
		graphicPath,
		graphicIfs,
	}
}

// BlurImage generate blur effect to an image.
func (graphic *Graphic) BlurImage(srcfile, dstfile string, sigma, numsteps float64, format string) (err error) {
	err = libgdkpixbuf.BlurImage(srcfile, dstfile, sigma, numsteps, libgdkpixbuf.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return
}

// ClipImage clip any recognized format image to target format image
// which could be "png" or "jpeg".
func (graphic *Graphic) ClipImage(srcfile, dstfile string, x, y, w, h int32, format string) (err error) {
	err = libgraphic.ClipImage(srcfile, dstfile, int(x), int(y), int(w), int(h), libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return
}

// ConvertImage converts from any recognized format imaget to target
// format image which could be "png" or "jpeg".
func (graphic *Graphic) ConvertImage(srcfile, dstfile, format string) (err error) {
	err = libgraphic.ConvertImage(srcfile, dstfile, libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return
}

// ConvertImageToDataUri convert image file to data uri.
func (graphic *Graphic) ConvertImageToDataUri(imgfile string) (dataUri string, err error) {
	return libgraphic.ConvertImageToDataUri(imgfile)
}

// ConvertDataUriToImage convert data uri to image file.
func (graphic *Graphic) ConvertDataUriToImage(dataUri string, dstfile string, format string) (err error) {
	return libgraphic.ConvertDataUriToImage(dataUri, dstfile, libgraphic.Format(format))
}

// CompositeImage composite two images.
func (graphic *Graphic) CompositeImage(srcfile, compfile, dstfile string, x, y int32, format string) (err error) {
	return libgraphic.CompositeImage(srcfile, compfile, dstfile, int(x), int(y), libgraphic.Format(format))
}

// CompositeImageUri composite two images which format in data uri.
func (graphic *Graphic) CompositeImageUri(srcDatauri, compDataUri string, x, y int32, format string) (dstDataUri string, err error) {
	return libgraphic.CompositeImageUri(srcDatauri, compDataUri, int(x), int(y), libgraphic.Format(format))
}

// GetDominantColorOfImage return the dominant hsv color of a image.
func (graphic *Graphic) GetDominantColorOfImage(imgfile string) (h, s, v float64, err error) {
	h, s, v, err = libgraphic.GetDominantColorOfImage(imgfile)
	if err != nil {
		logger.Errorf("%v", err)
	}
	return
}

// FillImage generate a new image in target width and height through
// source image, there are many fill sytles to choice from, such as
// "tile", "center", "stretch", "scalestretch".
func (graphic *Graphic) FillImage(srcfile, dstfile string, width, height int32, style, format string) (err error) {
	err = libgraphic.FillImage(srcfile, dstfile, int(width), int(height), libgraphic.FillStyle(style), libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return
}

// FlipImageHorizontal flip image in horizontal direction, and save as
// target format which could be "png" or "jpeg".
func (graphic *Graphic) FlipImageHorizontal(srcfile, dstfile string, format string) (err error) {
	err = libgraphic.FlipImageHorizontal(srcfile, dstfile, libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return
}

// FlipImageVertical flip image in vertical direction, and save as
// target format which could be "png" or "jpeg".
func (graphic *Graphic) FlipImageVertical(srcfile, dstfile string, format string) (err error) {
	err = libgraphic.FlipImageVertical(srcfile, dstfile, libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return
}

// Rgb2Hsv convert color format from RGB(r, g, b=[0..255]) to HSV(h=[0..360), s,v=[0..1]).
func (graphic *Graphic) Rgb2Hsv(r, g, b uint8) (h, s, v float64) {
	return libgraphic.Rgb2Hsv(r, g, b)
}

// Hsv2Rgb convert color format from HSV(h=[0..360), s,v=[0..1]) to RGB(r, g, b=[0..255]).
func (graphic *Graphic) Hsv2Rgb(h, s, v float64) (r, g, b uint8) {
	return libgraphic.Hsv2Rgb(h, s, v)
}

// GetImageSize return a image's width and height.
func (graphic *Graphic) GetImageSize(imgfile string) (int32, int32, error) {
	w, h, err := libgraphic.GetImageSize(imgfile)
	if err != nil {
		logger.Errorf("%v", err)
	}
	return int32(w), int32(h), err
}

// ResizeImage returns a new image file with the given width and
// height created by resizing the given image, and save to target
// image format which could be "png" or "jpeg".
func (graphic *Graphic) ResizeImage(srcfile, dstfile string, newWidth, newHeight int32, format string) (err error) {
	err = libgraphic.ScaleImage(srcfile, dstfile, int(newWidth), int(newHeight), libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return
}

// ThumbnailImage scale target image with limited maximum width and height.
func (graphic *Graphic) ThumbnailImage(srcfile, dstfile string, maxWidth, maxHeight uint32, format string) (err error) {
	err = libgraphic.ThumbnailImage(srcfile, dstfile, int(maxWidth), int(maxHeight), libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return
}

// RotateImageLeft rotate image to left side, and save to target image
// format which could be "png" or "jpeg".
func (graphic *Graphic) RotateImageLeft(srcfile, dstfile string, format string) (err error) {
	err = libgraphic.RotateImageLeft(srcfile, dstfile, libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return
}

// RotateImageRight rotate image to right side, and save to target image
// format which could be "png" or "jpeg".
func (graphic *Graphic) RotateImageRight(srcfile, dstfile string, format string) (err error) {
	err = libgraphic.RotateImageRight(srcfile, dstfile, libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return
}
