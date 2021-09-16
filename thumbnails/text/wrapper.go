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

// TODO: use 'code.google.com/p/freetype-go/freetype' to draw text on image
package text

// #cgo pkg-config: gdk-3.0
// #include <stdlib.h>
// #include "text.h"
import "C"

import (
	"bufio"
	"fmt"
	"os"
	"unsafe"

	"pkg.deepin.io/dde/api/thumbnails/loader"
)

type thumbInfo struct {
	width        int
	height       int
	xborder      int
	yborder      int
	canvasWidth  int
	canvasHeight int
	pixelsize    int
	fontsize     int
}

const (
	defaultDPI   int = 96
	defaultScale int = 1
)

// refer to kde text thumbnail in kde-runtime
func getThumbInfo(width, height int) *thumbInfo {
	var info thumbInfo
	// look good at width/height = 3/4
	if height*3 > width*4 {
		info.height = width * 4 / 3
		info.width = width
	} else {
		info.width = height * 3 / 4
		info.height = height
	}

	// one pixel for the rectangle, the rest. whitespace
	info.xborder = 1 + info.width/16  //minimum x-border
	info.yborder = 1 + info.height/16 //minimum y-border

	// calculate a better border so that the text is centered
	// kde: canvasWidth/canvasHeight = width/height - 2 * xborder / yborder
	info.canvasWidth = info.width - info.xborder
	info.canvasHeight = info.height - info.yborder

	// this font is supposed to look good at small sizes
	// pixelsize = Max(7, Min(10, (height - 2 * yborder) / 16))
	tmpPixel := (info.height - info.yborder) / 16
	if tmpPixel > 10 {
		tmpPixel = 10
	}
	info.pixelsize = tmpPixel
	// pixelsize = (fontsize * scale * dpi) / 72 from fontconfig
	info.fontsize = info.pixelsize * 72 / (defaultDPI * defaultScale)

	return &info
}
func doGenThumbnail(src, dest string, width, height int) (string, error) {
	info := getThumbInfo(width, height)
	strv, err := readFile(src, info)
	if err != nil {
		return "", err
	}
	defer freeCStrv(strv)

	tmp := loader.GetTmpImage()
	cTmp := C.CString(tmp)
	defer C.free(unsafe.Pointer(cTmp))

	var cinfo *C.ThumbInfo = &C.ThumbInfo{
		width:        C.int(info.width),
		height:       C.int(info.height),
		xborder:      C.int(info.xborder),
		yborder:      C.int(info.yborder),
		canvasWidth:  C.int(info.canvasWidth),
		canvasHeight: C.int(info.canvasHeight),
		fontSize:     C.int(info.fontsize),
	}
	ret := C.text_thumbnail(&strv[0], cTmp, cinfo)
	if ret != 0 {
		return "", fmt.Errorf("Draw text on image failed")
	}

	defer os.Remove(tmp)
	err = loader.ThumbnailImage(tmp, dest, width, height)
	if err != nil {
		return "", err
	}

	return dest, nil
}

func readFile(file string, info *thumbInfo) ([]*C.char, error) {
	fr, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fr.Close()
	}()

	var (
		cnt   int
		lines []string
	)
	scanner := bufio.NewScanner(fr)
	numLines := info.height / info.pixelsize
	for scanner.Scan() {
		if cnt >= numLines {
			break
		}
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		lines = append(lines, line)
		cnt += 1
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("Empty file")
	}

	bytesPerLine := info.width / info.pixelsize * 2
	return strvToCStrv(lines, bytesPerLine), nil
}

func strvToCStrv(strv []string, bytesPerLine int) []*C.char {
	var cstrv []*C.char
	for _, str := range strv {
		var tmp string
		for _, ch := range str {
			if len(tmp) > bytesPerLine {
				cstrv = append(cstrv, C.CString(tmp))
				tmp = ""
			}
			// convert 'tab' to 4 space
			if ch == '\t' {
				tmp += "    "
				continue
			}
			tmp += string(ch)
		}
		cstrv = append(cstrv, C.CString(tmp))
	}
	cstrv = append(cstrv, nil)
	return cstrv
}

func freeCStrv(strv []*C.char) {
	for _, str := range strv {
		C.free(unsafe.Pointer(str))
	}
}
