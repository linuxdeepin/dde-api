// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"github.com/godbus/dbus/v5"
	"github.com/linuxdeepin/go-lib/dbusutil"
	libgdkpixbuf "github.com/linuxdeepin/go-lib/gdkpixbuf"
	libgraphic "github.com/linuxdeepin/go-lib/graphic"
)

//go:generate dbusutil-gen em -type Graphic

const (
	dbusServiceName = "org.deepin.dde.Graphic1"
	dbusPath        = "/org/deepin/dde/Graphic1"
	dbusInterface   = "org.deepin.dde.Graphic1"
)

// Graphic is a dbus interface wrapper for github.com/linuxdeepin/go-lib/graphic.
type Graphic struct {
	service *dbusutil.Service
}

func (*Graphic) GetInterfaceName() string {
	return dbusInterface
}

// BlurImage generate blur effect to an image.
func (graphic *Graphic) BlurImage(srcFile, dstFile string, sigma, numSteps float64, format string) *dbus.Error {
	graphic.service.DelayAutoQuit()
	err := libgdkpixbuf.BlurImage(srcFile, dstFile, sigma, numSteps, libgdkpixbuf.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return dbusutil.ToError(err)
}

// ClipImage clip any recognized format image to target format image
// which could be "png" or "jpeg".
func (graphic *Graphic) ClipImage(srcFile, dstFile string, x, y, w, h int32, format string) *dbus.Error {
	graphic.service.DelayAutoQuit()
	err := libgraphic.ClipImage(srcFile, dstFile, int(x), int(y), int(w), int(h), libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return dbusutil.ToError(err)
}

// ConvertImage converts from any recognized format imaget to target
// format image which could be "png" or "jpeg".
func (graphic *Graphic) ConvertImage(srcFile, dstFile, format string) *dbus.Error {
	graphic.service.DelayAutoQuit()
	err := libgraphic.ConvertImage(srcFile, dstFile, libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return dbusutil.ToError(err)
}

// ConvertImageToDataUri convert image file to data uri.
func (graphic *Graphic) ConvertImageToDataUri(imgfile string) (dataUri string, busErr *dbus.Error) {
	graphic.service.DelayAutoQuit()
	dataUri, err := libgraphic.ConvertImageToDataUri(imgfile)
	return dataUri, dbusutil.ToError(err)
}

// ConvertDataUriToImage convert data uri to image file.
func (graphic *Graphic) ConvertDataUriToImage(dataUri string, dstFile string, format string) *dbus.Error {
	graphic.service.DelayAutoQuit()
	err := libgraphic.ConvertDataUriToImage(dataUri, dstFile, libgraphic.Format(format))
	return dbusutil.ToError(err)
}

// CompositeImage composite two images.
func (graphic *Graphic) CompositeImage(srcFile, compFile, dstFile string, x, y int32, format string) *dbus.Error {
	graphic.service.DelayAutoQuit()
	err := libgraphic.CompositeImage(srcFile, compFile, dstFile, int(x), int(y), libgraphic.Format(format))
	return dbusutil.ToError(err)
}

// CompositeImageUri composite two images which format in data uri.
func (graphic *Graphic) CompositeImageUri(srcDataUri, compDataUri string, x, y int32, format string) (resultDataUri string, busErr *dbus.Error) {
	graphic.service.DelayAutoQuit()
	resultDataUri, err := libgraphic.CompositeImageUri(srcDataUri, compDataUri, int(x), int(y), libgraphic.Format(format))
	return resultDataUri, dbusutil.ToError(err)
}

// GetDominantColorOfImage return the dominant hsv color of a image.
func (graphic *Graphic) GetDominantColorOfImage(imgFile string) (h, s, v float64, busErr *dbus.Error) {
	graphic.service.DelayAutoQuit()
	h, s, v, err := libgraphic.GetDominantColorOfImage(imgFile)
	if err != nil {
		logger.Errorf("%v", err)
	}
	return h, s, v, dbusutil.ToError(err)
}

// FillImage generate a new image in target width and height through
// source image, there are many fill sytles to choice from, such as
// "tile", "center".
func (graphic *Graphic) FillImage(srcFile, dstFile string,
	width, height int32, style, format string) *dbus.Error {

	graphic.service.DelayAutoQuit()
	err := libgraphic.FillImage(srcFile, dstFile, int(width), int(height), libgraphic.FillStyle(style), libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return dbusutil.ToError(err)
}

// FlipImageHorizontal flip image in horizontal direction, and save as
// target format which could be "png" or "jpeg".
func (graphic *Graphic) FlipImageHorizontal(srcFile, dstFile string, format string) *dbus.Error {
	graphic.service.DelayAutoQuit()
	err := libgraphic.FlipImageHorizontal(srcFile, dstFile, libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return dbusutil.ToError(err)
}

// FlipImageVertical flip image in vertical direction, and save as
// target format which could be "png" or "jpeg".
func (graphic *Graphic) FlipImageVertical(srcFile, dstFile string, format string) *dbus.Error {
	graphic.service.DelayAutoQuit()
	err := libgraphic.FlipImageVertical(srcFile, dstFile, libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return dbusutil.ToError(err)
}

// Rgb2Hsv convert color format from RGB(r, g, b=[0..255]) to HSV(h=[0..360), s,v=[0..1]).
func (graphic *Graphic) Rgb2Hsv(r, g, b uint8) (h, s, v float64, err *dbus.Error) {
	graphic.service.DelayAutoQuit()
	h, s, v = libgraphic.Rgb2Hsv(r, g, b)
	return
}

// Hsv2Rgb convert color format from HSV(h=[0..360), s,v=[0..1]) to RGB(r, g, b=[0..255]).
func (graphic *Graphic) Hsv2Rgb(h, s, v float64) (r, g, b uint8, err *dbus.Error) {
	graphic.service.DelayAutoQuit()
	r, g, b = libgraphic.Hsv2Rgb(h, s, v)
	return
}

// GetImageSize return a image's width and height.
func (graphic *Graphic) GetImageSize(imgFile string) (width int32, height int32, busErr *dbus.Error) {
	graphic.service.DelayAutoQuit()
	w, h, err := libgraphic.GetImageSize(imgFile)
	if err != nil {
		logger.Errorf("%v", err)
	}
	return int32(w), int32(h), dbusutil.ToError(err)
}

// ResizeImage returns a new image file with the given width and
// height created by resizing the given image, and save to target
// image format which could be "png" or "jpeg".
func (graphic *Graphic) ResizeImage(srcFile, dstFile string, newWidth, newHeight int32, format string) *dbus.Error {
	graphic.service.DelayAutoQuit()
	err := libgraphic.ScaleImage(srcFile, dstFile, int(newWidth), int(newHeight), libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return dbusutil.ToError(err)
}

// ThumbnailImage scale target image with limited maximum width and height.
func (graphic *Graphic) ThumbnailImage(srcFile, dstFile string,
	maxWidth, maxHeight uint32, format string) *dbus.Error {

	graphic.service.DelayAutoQuit()
	err := libgraphic.ThumbnailImage(srcFile, dstFile, int(maxWidth), int(maxHeight), libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return dbusutil.ToError(err)
}

// RotateImageLeft rotate image to left side, and save to target image
// format which could be "png" or "jpeg".
func (graphic *Graphic) RotateImageLeft(srcFile, dstFile string, format string) *dbus.Error {
	graphic.service.DelayAutoQuit()
	err := libgraphic.RotateImageLeft(srcFile, dstFile, libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return dbusutil.ToError(err)
}

// RotateImageRight rotate image to right side, and save to target image
// format which could be "png" or "jpeg".
func (graphic *Graphic) RotateImageRight(srcFile, dstFile string, format string) *dbus.Error {
	graphic.service.DelayAutoQuit()
	err := libgraphic.RotateImageRight(srcFile, dstFile, libgraphic.Format(format))
	if err != nil {
		logger.Errorf("%v", err)
	}
	return dbusutil.ToError(err)
}
