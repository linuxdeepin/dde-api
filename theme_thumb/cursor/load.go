// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package cursor

/*
#cgo pkg-config: xcursor
#cgo CFLAGS: -W -Wall -fPIC -fstack-protector-all
#include <stdlib.h>
#include <X11/Xcursor/Xcursor.h>
*/
import "C"
import (
	"fmt"
	"image"
	"image/color"
	"unsafe"

	"github.com/nfnt/resize"
)

func loadXCursorImage(filename string, size int) *C.XcursorImage {
	cFilename := C.CString(filename)
	xcImg := C.XcursorFilenameLoadImage(cFilename, C.int(size))
	C.free(unsafe.Pointer(cFilename))
	return xcImg
}

func destroyXCursorImage(img *C.XcursorImage) {
	C.XcursorImageDestroy(img)
}

func newImageFromXCurorImage(img *C.XcursorImage) image.Image {
	width := int(img.width)
	height := int(img.height)
	n := width * height
	// NOTE: (1 << 12) > 48*48
	pixels := (*[1 << 12]C.XcursorPixel)(unsafe.Pointer(img.pixels))[:n:n]
	newImg := image.NewRGBA(image.Rect(0, 0, int(img.width), int(img.height)))
	var i = 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := pixels[i]
			// pixel format: ARGB
			alpha := uint8(pixel >> 24)
			red := uint8((pixel >> 16) & 0xff)
			green := uint8((pixel >> 8) & 0xff)
			blue := uint8(pixel & 0xff)
			color := color.RGBA{R: red, G: green, B: blue, A: alpha}
			newImg.SetRGBA(x, y, color)
			i++
		}
	}
	return newImg
}

func loadXCursor(filename string, size int) (image.Image, error) {
	xcImg := loadXCursorImage(filename, size)
	if xcImg == nil {
		return nil, fmt.Errorf("failed to load x cursor image %q", filename)
	}
	defer destroyXCursorImage(xcImg)
	img := newImageFromXCurorImage(xcImg)
	imgWidth := img.Bounds().Dx()
	if imgWidth != size {
		img = resize.Resize(uint(size), 0, img, resize.Bilinear)
	}
	return img, nil
}
